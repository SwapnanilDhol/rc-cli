package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initWebhooks() {
	webhooksCmd := &cobra.Command{
		Use:   "webhooks",
		Short: "Manage webhooks",
	}

	webhooksListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		RunE:  runListWebhooks,
	}

	webhooksCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new webhook",
		RunE:  runCreateWebhook,
	}

	webhooksTestCmd := &cobra.Command{
		Use:   "test",
		Short: "Send a test webhook",
		RunE:  runTestWebhook,
	}

	webhooksEventsCmd := &cobra.Command{
		Use:   "events",
		Short: "List webhook events",
		RunE:  runWebhookEvents,
	}

	webhooksCmd.AddCommand(webhooksListCmd, webhooksCreateCmd, webhooksTestCmd, webhooksEventsCmd)
	RootCmd.AddCommand(webhooksCmd)

	// Aliases
	RootCmd.AddCommand(&cobra.Command{
		Use:   "webhooks:list",
		Short: "List all webhooks",
		RunE:  runListWebhooks,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use:   "webhooks:create",
		Short: "Create a new webhook",
		RunE:  runCreateWebhook,
	})
}

func runListWebhooks(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n🔗 Fetching webhooks...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/webhooks", cfg.ProjectID))
	if err != nil {
		fmt.Println(appsStyle.Render("\n🔗 Webhooks\n"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Webhooks are managed in the RevenueCat dashboard"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/webhooks\n"))
		return nil
	}

	fmt.Println(appsStyle.Render("\n🔗 Webhooks\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			webhook := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", webhook["id"])))
			fmt.Printf("  URL: %s\n", webhook["url"])
			fmt.Printf("  Events: %v\n", webhook["events"])
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  No webhooks configured."))
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/webhooks to create one\n"))
	}

	return nil
}

func runCreateWebhook(cmd *cobra.Command, args []string) error {
	fmt.Println(appsStyle.Render("\n🔗 Create Webhook\n"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Creating webhooks is done in the RevenueCat dashboard"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/webhooks\n"))

	return nil
}

func runTestWebhook(cmd *cobra.Command, args []string) error {
	fmt.Println(appsStyle.Render("\n🔗 Test Webhook\n"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Testing webhooks is done in the RevenueCat dashboard"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/webhooks\n"))

	return nil
}

func runWebhookEvents(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Fetching webhook events...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/webhooks/events", cfg.ProjectID))
	if err != nil {
		fmt.Println(appsStyle.Render("\n📋 Webhook Events\n"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Viewing webhook events is done in the RevenueCat dashboard"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/webhooks/events\n"))
		return nil
	}

	fmt.Println(appsStyle.Render("\n📋 Webhook Events\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			event := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", event["id"])))
			fmt.Printf("  Type: %s\n", event["type"])
			fmt.Printf("  Status: %s\n", event["status"])
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  No webhook events found.\n"))
	}

	return nil
}
