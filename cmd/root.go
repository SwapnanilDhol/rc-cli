package cmd

import (
	"github.com/spf13/cobra"
)

var (
	apiKey    string
	projectID string

	RootCmd = &cobra.Command{
		Use:   "rc",
		Short: "RevenueCat CLI - Manage your subscription apps, users, and analytics",
		Long:  `RevenueCat CLI for managing subscriptions, apps, and analytics.`,
	}
)

func Execute() error {
	return RootCmd.Execute()
}

var internalCmd = &cobra.Command{
	Use:   "internal",
	Short: "Internal dashboard API commands (app.revenuecat.com/internal/v1)",
	Long:  `Commands for the internal RevenueCat dashboard API. These commands use session-based authentication and are for managing projects, offerings, entitlements, and analytics through the dashboard backend.`,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "RevenueCat API key")
	RootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "p", "", "RevenueCat project ID")

	RootCmd.AddCommand(internalCmd)

	// Initialize all subcommands
	initApps()
	initConfig()
	initProjects()
	initSubscribers()
	initProducts()
	initSubscriptions()
	initEntitlements()
	initWebhooks()
	initApiV2()
	// initInternal() is called automatically via its own init() in cmd/internal.go
}
