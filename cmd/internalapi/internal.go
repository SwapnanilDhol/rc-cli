package internalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/config"
	rcinternal "revenuecat-cli/internal"
)

var (
	InternalStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	GreenStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	YellowStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	GrayStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	CyanStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

// Options are the global flags. They travel through the command context rather
// than package-level variables, so nothing here is mutable shared state and
// tests do not depend on execution order.
type Options struct {
	ProjectRef string // -p/--project-id: an ID or a name
	JSON       bool   // --json
}

type optionsKey struct{}

// WithOptions attaches the global flags to a command context. Called once, by
// the root command's PersistentPreRun.
func WithOptions(ctx context.Context, o Options) context.Context {
	return context.WithValue(ctx, optionsKey{}, o)
}

// OptionsOf reads the global flags back out. A zero Options is a valid default,
// so commands invoked outside the root (in tests) still work.
func OptionsOf(cmd *cobra.Command) Options {
	if cmd == nil || cmd.Context() == nil {
		return Options{}
	}
	o, _ := cmd.Context().Value(optionsKey{}).(Options)
	return o
}

// Exported command groups for the internal CLI shim
var (
	ProjectsCmd      = projectsCmd
	EntitlementsCmd  = entitlementsCmd
	OfferingsCmd     = offeringsCmd
	ProductsCmd      = productsCmd
	AppsCmd          = appsCmd
	ChartsCmd        = chartsCmd
	ListsCmd         = listsCmd
	ExperimentsCmd   = experimentsCmd
	StoresStatusCmd  = storesStatusCmd
	CollaboratorsCmd = collaboratorsCmd
	APIKeysCmd       = apiKeysCmd
	AuditCmd         = auditCmd
	UtilitiesCmd     = utilitiesCmd
)

func GetInternalClient() (*rcinternal.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	client, refreshed, err := rcinternal.EnsureAuthenticated(cfg)
	if err != nil {
		return nil, err
	}
	if refreshed {
		if saveErr := config.SaveConfig(cfg); saveErr != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", saveErr)
		}
	}
	return client, nil
}

// GetProjectID resolves the project for this invocation: the -p flag if given,
// otherwise the saved default. A name is resolved to an ID.
func GetProjectID(cmd *cobra.Command) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	ref := OptionsOf(cmd).ProjectRef
	if ref == "" {
		ref = cfg.ProjectID
	}
	if ref == "" {
		// Report the missing credential first — telling someone to run
		// 'projects list' when they are not logged in just fails again one
		// step later.
		if cfg.AuthToken == "" && (cfg.Email == "" || cfg.Password == "") {
			return "", fmt.Errorf("not logged in. Run: rc login")
		}
		return "", fmt.Errorf("no project selected. Run: rc internal projects list && rc internal projects use -i <project_id>")
	}
	if looksLikeID(ref) {
		return ref, nil
	}
	// -p was given a project name rather than an ID; resolve it.
	client, err := GetInternalClient()
	if err != nil {
		return "", err
	}
	return ResolveProjectID(client, ref)
}

func ToJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func GetStringValue(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return fallback
}

// Progress prints a status line to stderr, and nothing at all under --json, so
// that stdout stays a clean parseable document.
func Progress(cmd *cobra.Command, msg string) {
	if OptionsOf(cmd).JSON {
		return
	}
	fmt.Fprintln(os.Stderr, msg)
}

// CheckResponse turns a non-2xx status or an API-level error code into a Go error.
// Without this a 401/403 decodes into an empty Response and prints as "no results".
func CheckResponse(resp *rcinternal.Response) error {
	if resp == nil {
		return fmt.Errorf("empty response from internal API")
	}
	if resp.StatusCode >= 400 {
		msg := resp.Message
		if msg == "" {
			msg = "request failed"
		}
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return fmt.Errorf("HTTP %d: %s (run: rc login)", resp.StatusCode, msg)
		}
		return fmt.Errorf("HTTP %d: %s %s", resp.StatusCode, resp.Code, msg)
	}
	if resp.Code != "" {
		return fmt.Errorf("%s: %s", resp.Code, resp.Message)
	}
	return nil
}

// Ctx is what every dashboard command needs: the resolved project and an
// authenticated client. Building it once removes the project/auth preamble that
// was otherwise repeated in every runner.
type Ctx struct {
	ProjectID string
	Client    *rcinternal.Client
	JSON      bool
}

// Dashboard resolves the project, authenticates, and prints the progress line.
// Project-scoped commands start here.
func Dashboard(cmd *cobra.Command, progress string) (*Ctx, error) {
	projectID, err := GetProjectID(cmd)
	if err != nil {
		return nil, err
	}
	client, err := GetInternalClient()
	if err != nil {
		return nil, err
	}
	Progress(cmd, progress)
	return &Ctx{ProjectID: projectID, Client: client, JSON: OptionsOf(cmd).JSON}, nil
}

// DashboardNoProject is for the handful of commands that are account-scoped
// rather than project-scoped (listing or creating projects, for instance).
func DashboardNoProject(cmd *cobra.Command, progress string) (*Ctx, error) {
	client, err := GetInternalClient()
	if err != nil {
		return nil, err
	}
	Progress(cmd, progress)
	return &Ctx{Client: client, JSON: OptionsOf(cmd).JSON}, nil
}

// Path builds a project-scoped internal API path.
func (c *Ctx) Path(format string, args ...interface{}) string {
	return fmt.Sprintf("/developers/me/projects/"+c.ProjectID+format, args...)
}

// Respond is the single place the output mode is decided, and the only place a
// command should return from after a request. It checks the response, then
// renders either JSON or the human view — never both, and never one instead of
// a side effect, because anything that must happen regardless belongs before
// this call.
func (c *Ctx) Respond(resp *rcinternal.Response, human func() error) error {
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if c.JSON {
		EmitJSONValue(resp.Payload())
		return nil
	}
	if human == nil {
		return nil
	}
	return human()
}

// EmitJSONValue writes a value to stdout as indented JSON.
func EmitJSONValue(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
	}
}

// hexID matches the bare-hex IDs the dashboard uses for projects (e.g. "a1b2c3d4").
// Prefixed IDs (ofrng…, prod…, app…) are recognised separately by looksLikeID.
var hexID = regexp.MustCompile(`^[0-9a-f]{6,}$`)

// idPrefixes are the known opaque-ID prefixes across dashboard entities. Note the
// dashboard is not consistent: offerings are "ofrng…" but projects are bare hex.
var idPrefixes = []string{"ofrng", "prod", "proj", "app", "entl", "pkg", "exp"}

// looksLikeID reports whether ref is already an ID, letting callers skip a list
// round trip. A false negative only costs one extra request; a false positive is
// caught because resolveByRef matches on id as well as name.
func looksLikeID(ref string) bool {
	if hexID.MatchString(ref) {
		return true
	}
	for _, p := range idPrefixes {
		if strings.HasPrefix(ref, p) && len(ref) > len(p) {
			return true
		}
	}
	return false
}

// resolveByRef finds an entity ID from a literal ID, an identifier, or a display
// name. It matches on id first so a literal ID always resolves regardless of the
// entity's ID format. Name matching is case-insensitive, and ambiguity is an error
// rather than a silent first-match, so a reference always maps to the same entity
// or fails loudly.
func resolveByRef(items []interface{}, ref, kind string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("%s reference is required", kind)
	}

	needle := strings.ToLower(ref)
	var matches []string
	var names []string
	for _, it := range items {
		obj, ok := it.(map[string]interface{})
		if !ok {
			continue
		}
		id := GetStringValue(obj, "id", "")
		if id == "" {
			continue
		}
		// An exact ID match wins outright — no ambiguity check needed.
		if id == ref {
			return id, nil
		}
		identifier := GetStringValue(obj, "identifier", "")
		displayName := GetStringValue(obj, "display_name", "")
		if displayName == "" {
			displayName = GetStringValue(obj, "name", "")
		}
		label := identifier
		if label == "" {
			label = displayName
		}
		names = append(names, fmt.Sprintf("%s (%s)", label, id))
		if strings.EqualFold(identifier, needle) || strings.EqualFold(displayName, needle) {
			matches = append(matches, id)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return "", fmt.Errorf("no %s matches %q. Available: %s", kind, ref, strings.Join(names, ", "))
	default:
		return "", fmt.Errorf("%q is ambiguous — it matches %d %ss (%s). Pass the %s ID instead",
			ref, len(matches), kind, strings.Join(matches, ", "), kind)
	}
}

// ResolveOfferingID accepts an offering ID (ofrng…) or an offering identifier /
// display name and returns the ID.
func ResolveOfferingID(client *rcinternal.Client, projectID, ref string) (string, error) {
	if strings.HasPrefix(ref, "ofrng") {
		return ref, nil
	}
	resp, err := client.Get(fmt.Sprintf("/developers/me/projects/%s/offerings", projectID))
	if err != nil {
		return "", err
	}
	if err := CheckResponse(resp); err != nil {
		return "", err
	}
	return resolveByRef(resp.Items, ref, "offering")
}

// ResolveProjectID accepts a project ID (proj…) or a project name and returns the ID.
func ResolveProjectID(client *rcinternal.Client, ref string) (string, error) {
	resp, err := client.Get("/developers/me/projects")
	if err != nil {
		return "", err
	}
	if err := CheckResponse(resp); err != nil {
		return "", err
	}
	return resolveByRef(resp.Items, ref, "project")
}

// ResolveAppID accepts an app ID (app…) or an app name and returns the ID.
func ResolveAppID(client *rcinternal.Client, projectID, ref string) (string, error) {
	resp, err := client.Get(fmt.Sprintf("/developers/me/projects/%s/apps", projectID))
	if err != nil {
		return "", err
	}
	if err := CheckResponse(resp); err != nil {
		return "", err
	}
	return resolveByRef(resp.Items, ref, "app")
}
