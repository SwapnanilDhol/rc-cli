package cmd

import (
	"github.com/spf13/cobra"
	rcinternalapi "revenuecat-cli/cmd/internalapi"
)

// initInternal registers the internal command group under the root.
// All internal commands use the internal dashboard API at https://app.revenuecat.com/internal/v1
// vs the public /v2 API. Internal commands require session cookie auth (rc login),
// while public API commands use API key auth (rc config).
func initInternal() {
	internalCmd := &cobra.Command{
		Use:   "internal",
		Short: "Internal Dashboard API commands",
		Long: `Commands for the Internal Dashboard API at https://app.revenuecat.com/internal/v1

This is separate from the public /v2 API. The internal API uses session cookie
authentication (rc login) vs the public API's API key authentication (rc config).

Available command groups:
  projects        - Manage projects
  entitlements    - Manage entitlements
  offerings       - Manage offerings
  products        - Manage products
  apps            - Manage apps
  charts          - View analytics charts
  lists           - Manage subscriber lists
  experiments     - Manage price experiments (A/B tests)
  stores-status   - Get product stores connection status
  collaborators   - List project collaborators
  apikeys         - List API keys
  audit           - View audit logs
  utilities       - Utility endpoints (countries)
  api             - Call any internal endpoint directly

All internal commands require authentication via 'rc login'.`,
	}

	// Register all internal subcommand groups
	internalCmd.AddCommand(
		rcinternalapi.ProjectsCmd,
		rcinternalapi.EntitlementsCmd,
		rcinternalapi.OfferingsCmd,
		rcinternalapi.ProductsCmd,
		rcinternalapi.AppsCmd,
		rcinternalapi.ChartsCmd,
		rcinternalapi.ListsCmd,
		rcinternalapi.ExperimentsCmd,
		rcinternalapi.StoresStatusCmd,
		rcinternalapi.CollaboratorsCmd,
		rcinternalapi.APIKeysCmd,
		rcinternalapi.AuditCmd,
		rcinternalapi.UtilitiesCmd,
		rcinternalapi.RawAPICmd,
	)

	RootCmd.AddCommand(internalCmd)
}
