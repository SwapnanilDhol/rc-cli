package cmd

import (
	"context"
	"io/fs"

	"github.com/spf13/cobra"
	"revenuecat-cli/cmd/internalapi"
)

const version = "1.1.0"

// NewRootCmd builds the whole command tree. Building it in a function rather than
// package init() means tests get a fresh, isolated tree and nothing depends on
// initialisation order.
//
// skills is the embedded agent-skill filesystem, supplied by the root package
// because go:embed cannot reach paths above its own directory. It may be nil, in
// which case `rc install-skills` reports that none are bundled.
func NewRootCmd(skills fs.FS, skillsRoot string) *cobra.Command {
	var opts options

	root := &cobra.Command{
		Use:     "rc",
		Short:   "RevenueCat CLI - manage subscriptions, offerings, and analytics",
		Version: version,
		Long: `RevenueCat CLI.

Two APIs, two credentials, one binary — which one a command uses is always visible
from the command itself:

  rc <command>            Public v2 API   api.revenuecat.com/v2
                          Auth: API key   → rc config
  rc internal <command>   Dashboard API   app.revenuecat.com/internal/v1
                          Auth: session   → rc login

Anything the dashboard can do but the public API cannot — offerings metadata,
packages, experiments, charts — lives under 'rc internal'.

Pass --json to any command for machine-readable output.`,
	}

	root.PersistentFlags().StringVar(&opts.APIKey, "api-key", "", "RevenueCat public v2 API key (overrides rc config)")
	root.PersistentFlags().StringVarP(&opts.ProjectID, "project-id", "p", "", "RevenueCat project ID or name (overrides saved default)")
	root.PersistentFlags().BoolVar(&opts.JSON, "json", false, "Emit machine-readable JSON instead of formatted text")

	// Publish the global flags on the command context, so neither package needs
	// mutable shared state to read them.
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx = withOptions(ctx, opts)
		ctx = internalapi.WithOptions(ctx, internalapi.Options{
			ProjectRef: opts.ProjectID,
			JSON:       opts.JSON,
		})
		cmd.SetContext(ctx)
	}

	// Public v2 API commands
	initApps(root)
	initConfig(root)
	initProjects(root)
	initSubscribers(root)
	initProducts(root)
	initSubscriptions(root)
	initEntitlements(root)
	initApiV2(root)
	initInstallSkills(root, skills, skillsRoot)

	// Dashboard API commands (rc internal …) plus rc login / rc logout
	initInternal(root)
	initAuth(root)

	return root
}

// Execute builds and runs the CLI.
func Execute(skills fs.FS, skillsRoot string) error {
	return NewRootCmd(skills, skillsRoot).Execute()
}
