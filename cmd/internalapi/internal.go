package internalapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"revenuecat-cli/config"
	rcinternal "revenuecat-cli/internal"
)

var (
	InternalStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	GreenStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	YellowStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	GrayStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	CyanStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	flagProjectID string // set via flags by root.go

	// JSONOutput makes every internal command emit machine-readable JSON on stdout
	// instead of the decorated human view. Set from the global --json flag.
	JSONOutput bool
)

// SetFlagProjectID is called by root.go to pass the -p/--project-id flag value
func SetFlagProjectID(id string) {
	flagProjectID = id
}

// SetJSONOutput is called by root.go to pass the --json flag value
func SetJSONOutput(v bool) {
	JSONOutput = v
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

func GetProjectID() (string, error) {
	// Flag (-p/--project-id) takes precedence over saved config
	ref := flagProjectID
	if ref == "" {
		cfg, err := config.LoadConfig()
		if err != nil {
			return "", err
		}
		ref = cfg.ProjectID
	}
	if ref == "" {
		return "", fmt.Errorf("no project selected. Run: rc internal projects list && rc internal projects use -i <project_id>")
	}
	if strings.HasPrefix(ref, "proj") {
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
func Progress(msg string) {
	if JSONOutput {
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

// EmitJSON prints the response payload as JSON when --json is set and reports
// whether it did, so callers can `if EmitJSON(resp) { return nil }` before their
// human-readable printing.
func EmitJSON(resp *rcinternal.Response) bool {
	if !JSONOutput {
		return false
	}
	var payload interface{}
	switch {
	case len(resp.Items) > 0:
		payload = resp.Items
	case resp.Data != nil:
		payload = resp.Data
	default:
		payload = map[string]interface{}{"status": resp.StatusCode}
	}
	return EmitJSONValue(payload)
}

// EmitJSONValue prints an arbitrary value as JSON when --json is set.
func EmitJSONValue(v interface{}) bool {
	if !JSONOutput {
		return false
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
	}
	return true
}

// resolveByRef finds an entity ID from either a literal ID (recognised by its
// prefix) or a human-facing identifier / display name. Name matching is
// case-insensitive and ambiguity is an error rather than a silent first-match,
// so a given reference always maps to the same entity or fails loudly.
func resolveByRef(items []interface{}, ref, idPrefix, kind string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("%s reference is required", kind)
	}
	if strings.HasPrefix(ref, idPrefix) {
		return ref, nil
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
	return resolveByRef(resp.Items, ref, "ofrng", "offering")
}

// ResolveProjectID accepts a project ID (proj…) or a project name and returns the ID.
func ResolveProjectID(client *rcinternal.Client, ref string) (string, error) {
	if strings.HasPrefix(ref, "proj") {
		return ref, nil
	}
	resp, err := client.Get("/developers/me/projects")
	if err != nil {
		return "", err
	}
	if err := CheckResponse(resp); err != nil {
		return "", err
	}
	return resolveByRef(resp.Items, ref, "proj", "project")
}

// ResolveAppID accepts an app ID (app…) or an app name and returns the ID.
func ResolveAppID(client *rcinternal.Client, projectID, ref string) (string, error) {
	if strings.HasPrefix(ref, "app") {
		return ref, nil
	}
	resp, err := client.Get(fmt.Sprintf("/developers/me/projects/%s/apps", projectID))
	if err != nil {
		return "", err
	}
	if err := CheckResponse(resp); err != nil {
		return "", err
	}
	return resolveByRef(resp.Items, ref, "app", "app")
}
