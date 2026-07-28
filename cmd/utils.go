package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"revenuecat-cli/config"
)

// options are the global flags for public v2 commands. Like the dashboard side,
// they travel on the command context instead of package-level variables.
type options struct {
	APIKey    string
	ProjectID string
	JSON      bool
}

type optionsKey struct{}

func withOptions(ctx context.Context, o options) context.Context {
	return context.WithValue(ctx, optionsKey{}, o)
}

func optionsOf(cmd *cobra.Command) options {
	if cmd == nil || cmd.Context() == nil {
		return options{}
	}
	o, _ := cmd.Context().Value(optionsKey{}).(options)
	return o
}

// jsonRequested reports whether --json was passed.
func jsonRequested(cmd *cobra.Command) bool { return optionsOf(cmd).JSON }

// emitJSON writes v to stdout as indented JSON. Used by every command when --json
// is set so agents get a parseable payload instead of the decorated human view.
func emitJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// loadConfig reads the saved config and applies this invocation's flag overrides.
func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	o := optionsOf(cmd)
	if o.APIKey != "" {
		cfg.APIKey = o.APIKey
	}
	if o.ProjectID != "" {
		cfg.ProjectID = o.ProjectID
	}
	return cfg, nil
}

func getStringValue(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return fallback
}

func getBoolValue(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func formatDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return "N/A"
}

func formatDateTime(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return "N/A"
}

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return home
}

func getConfigPath() string {
	return fmt.Sprintf("%s/.revenuerc", userHomeDir())
}

func formatUnixTime(ts int64) string {
	if ts == 0 {
		return "N/A"
	}
	return time.Unix(ts/1000, 0).Format("2006-01-02")
}
