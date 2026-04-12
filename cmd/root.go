package cmd

import (
	"github.com/spf13/cobra"
)

var (
	apiKey    string
	projectID string

	version = "0.2.0"

	RootCmd = &cobra.Command{
		Use:   "rc",
		Short: "RevenueCat CLI - Manage your subscription apps, users, and analytics",
		Long:  `RevenueCat CLI for managing subscriptions, apps, and analytics.`,
		Version: version,
	}
)

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "RevenueCat API key")
	RootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "p", "", "RevenueCat project ID")

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
}
