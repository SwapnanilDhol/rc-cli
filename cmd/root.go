package cmd

import (
	"github.com/spf13/cobra"
	"revenuecat-cli/cmd/internalapi"
)

var (
	apiKey     string
	projectID  string
	jsonOutput bool

	version = "0.3.0"

	RootCmd = &cobra.Command{
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
)

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "RevenueCat public v2 API key (overrides rc config)")
	RootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "p", "", "RevenueCat project ID (overrides saved default)")
	RootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON instead of formatted text")

	// Hand the global flags to the internal command tree, which lives in its own
	// package and cannot see these vars directly.
	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		internalapi.SetFlagProjectID(projectID)
		internalapi.SetJSONOutput(jsonOutput)
	}

	// Public v2 API commands
	initApps()
	initConfig()
	initProjects()
	initSubscribers()
	initProducts()
	initSubscriptions()
	initEntitlements()
	initOffers()
	initWebhooks()
	initApiV2()
	initInstallSkills()

	// Dashboard API commands (rc internal …) plus rc login / rc logout
	initInternal()
}
