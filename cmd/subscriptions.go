package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initSubscriptions() {
	subscriptionsCmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Manage subscriptions",
	}

	subscriptionsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List subscriptions for a customer",
		RunE:  runListSubscriptions,
	}
	subscriptionsListCmd.Flags().StringP("app-user-id", "u", "", "App User ID")

	subscriptionsCmd.AddCommand(subscriptionsListCmd)
	RootCmd.AddCommand(subscriptionsCmd)

	// Alias
	RootCmd.AddCommand(&cobra.Command{
		Use:   "subs:list",
		Short: "List subscriptions",
		RunE:  runListSubscriptions,
	})
}

func runListSubscriptions(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	appUserID, _ := cmd.Flags().GetString("app-user-id")
	if appUserID == "" {
		return fmt.Errorf("app-user-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Fetching subscriptions...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/customers/%s", cfg.ProjectID, appUserID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render(fmt.Sprintf("\n📦 Subscriptions for %s:\n", appUserID)))

	if data, ok := resp.Data.(map[string]interface{}); ok {
		if subscriptions, ok := data["subscriptions"].(map[string]interface{}); ok {
			if items, ok := subscriptions["items"].([]interface{}); ok && len(items) > 0 {
				for _, item := range items {
					sub := item.(map[string]interface{})
					fmt.Println(cyanStyle.Render(fmt.Sprintf("  Product: %s", sub["product_id"])))
					fmt.Printf("    Store: %s\n", sub["store"])
					fmt.Printf("    Status: %s\n", sub["status"])
					fmt.Println()
				}
			}
		}

		if entitlements, ok := data["active_entitlements"].(map[string]interface{}); ok {
			if items, ok := entitlements["items"].([]interface{}); ok && len(items) > 0 {
				hasActiveSubs := false
				if subscriptions, ok := data["subscriptions"].(map[string]interface{}); ok {
					if subsItems, ok := subscriptions["items"].([]interface{}); ok && len(subsItems) > 0 {
						hasActiveSubs = true
					}
				}

				if !hasActiveSubs {
					fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  No active subscriptions, but has entitlements.\n"))
				}

				for _, item := range items {
					ent := item.(map[string]interface{})
					fmt.Println(cyanStyle.Render(fmt.Sprintf("  Entitlement: %s", ent["identifier"])))
					fmt.Printf("    Product: %s\n", ent["product_identifier"])
					fmt.Println()
				}
			}
		}
	} else if len(resp.Items) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  No subscriptions found."))
	}

	return nil
}
