package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initSubscribers() {
	subscribersCmd := &cobra.Command{
		Use:   "subscribers",
		Short: "Manage subscribers (customers)",
	}

	subscribersListCmd := &cobra.Command{
		Use:   "list",
		Short: "List subscribers",
		RunE:  runListSubscribers,
	}
	subscribersListCmd.Flags().StringP("limit", "l", "20", "Number of results (max 200)")
	subscribersListCmd.Flags().StringP("platform", "p", "", "Filter by platform (ios, android, stripe)")
	subscribersListCmd.Flags().StringP("product-id", "r", "", "Filter by product ID")
	subscribersListCmd.Flags().StringP("starting-after", "s", "", "Cursor for pagination (customer ID)")

	subscribersGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get subscriber details by app_user_id",
		RunE:  runGetSubscriber,
	}
	subscribersGetCmd.Flags().StringP("app-user-id", "u", "", "App User ID (or $RCAnonymousID:...)")

	subscribersSearchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for a subscriber",
		RunE:  runSearchSubscriber,
	}
	subscribersSearchCmd.Flags().StringP("query", "q", "", "Search query")

	subscribersEntitlementsCmd := &cobra.Command{
		Use:   "entitlements",
		Short: "Get subscriber's active entitlements",
		RunE:  runSubscriberEntitlements,
	}
	subscribersEntitlementsCmd.Flags().StringP("app-user-id", "u", "", "App User ID")

	subscribersSubscriptionsCmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Get subscriber's subscriptions",
		RunE:  runSubscriberSubscriptions,
	}
	subscribersSubscriptionsCmd.Flags().StringP("app-user-id", "u", "", "App User ID")

	subscribersCmd.AddCommand(subscribersListCmd, subscribersGetCmd, subscribersSearchCmd, subscribersEntitlementsCmd, subscribersSubscriptionsCmd)
	RootCmd.AddCommand(subscribersCmd)

	// Aliases
	RootCmd.AddCommand(&cobra.Command{
		Use:   "users:list",
		Short: "List subscribers",
		RunE:  runListSubscribers,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use:   "users:get",
		Short: "Get subscriber details",
		RunE:  runGetSubscriber,
	})
}

func runListSubscribers(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	limit, _ := cmd.Flags().GetString("limit")
	platform, _ := cmd.Flags().GetString("platform")
	productID, _ := cmd.Flags().GetString("product-id")
	startingAfter, _ := cmd.Flags().GetString("starting-after")

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n👥 Fetching subscribers...")

	params := map[string]string{"limit": limit}
	if platform != "" {
		params["platform"] = strings.ToLower(platform)
	}
	if productID != "" {
		params["products"] = productID
	}
	if startingAfter != "" {
		params["starting_after"] = startingAfter
	}

	path := fmt.Sprintf("/projects/%s/customers", cfg.ProjectID)
	resp, err := client.GetWithParams(path, params)
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n👥 Subscribers:\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			sub := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", sub["id"])))
			fmt.Printf("  Platform: %v\n", getStringValue(sub, "last_seen_platform", "N/A"))
			fmt.Printf("  App Version: %v\n", getStringValue(sub, "last_seen_app_version", "N/A"))
			fmt.Printf("  Country: %v\n", getStringValue(sub, "last_seen_country", "N/A"))
			fmt.Printf("  First Seen: %v\n", formatDate(getStringValue(sub, "first_seen_at", "")))
			fmt.Printf("  Last Seen: %v\n", formatDate(getStringValue(sub, "last_seen_at", "")))
			fmt.Println()
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("Showing %d subscribers", len(resp.Items))))
		if resp.NextPage != "" {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  More results available. Use --starting-after <last-id> for next page"))
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No subscribers found."))
	}

	return nil
}

func runGetSubscriber(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n👤 Fetching subscriber details...")

	// URL encode the customer ID
	encodedID := url.QueryEscape(appUserID)
	path := fmt.Sprintf("/projects/%s/customers/%s", cfg.ProjectID, encodedID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		fmt.Println(appsStyle.Render("\n👤 Subscriber Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", data["id"])))
		fmt.Printf("  Platform: %v\n", getStringValue(data, "last_seen_platform", "N/A"))
		fmt.Printf("  App Version: %v\n", getStringValue(data, "last_seen_app_version", "N/A"))
		fmt.Printf("  Country: %v\n", getStringValue(data, "last_seen_country", "N/A"))
		fmt.Printf("  First Seen: %v\n", formatDateTime(getStringValue(data, "first_seen_at", "")))
		fmt.Printf("  Last Seen: %v\n", formatDateTime(getStringValue(data, "last_seen_at", "")))

		if entitlements, ok := data["active_entitlements"].(map[string]interface{}); ok {
			if items, ok := entitlements["items"].([]interface{}); ok && len(items) > 0 {
				fmt.Println(appsStyle.Render("\n📋 Active Entitlements:"))
				for _, item := range items {
					ent := item.(map[string]interface{})
					fmt.Printf("\n  %s\n", cyanStyle.Render(getStringValue(ent, "identifier", "N/A")))
					fmt.Printf("    Product ID: %s\n", getStringValue(ent, "product_identifier", "N/A"))
					fmt.Printf("    Expires: %s\n", formatDateTime(getStringValue(ent, "expires_at", "")))
					fmt.Printf("    Period Type: %s\n", getStringValue(ent, "period_type", "N/A"))
				}
			}
		}
	}

	return nil
}

func runSearchSubscriber(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	query, _ := cmd.Flags().GetString("query")
	if query == "" {
		return fmt.Errorf("query is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n🔍 Searching for subscriber...")

	params := map[string]string{"limit": "200"}
	resp, err := client.GetWithParams(fmt.Sprintf("/projects/%s/customers", cfg.ProjectID), params)
	if err != nil {
		return err
	}

	var matches []map[string]interface{}
	queryLower := strings.ToLower(query)
	for _, item := range resp.Items {
		sub := item.(map[string]interface{})
		if id, ok := sub["id"].(string); ok {
			if strings.Contains(strings.ToLower(id), queryLower) {
				matches = append(matches, sub)
			}
		}
	}

	if len(matches) > 0 {
		fmt.Println(appsStyle.Render(fmt.Sprintf("\n🔍 Search Results (%d):\n", len(matches))))
		for _, sub := range matches {
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", sub["id"])))
			fmt.Printf("  Platform: %v\n", getStringValue(sub, "last_seen_platform", "N/A"))
			fmt.Printf("  Last Seen: %v\n", formatDate(getStringValue(sub, "last_seen_at", "")))
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No subscriber found."))
	}

	return nil
}

func runSubscriberEntitlements(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n📋 Fetching entitlements...")

	encodedID := url.QueryEscape(appUserID)
	path := fmt.Sprintf("/projects/%s/customers/%s/active_entitlements", cfg.ProjectID, encodedID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render(fmt.Sprintf("\n📋 Entitlements for %s:\n", appUserID)))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			ent := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  %s", ent["identifier"])))
			fmt.Printf("    Product: %s\n", getStringValue(ent, "product_identifier", "N/A"))
			fmt.Printf("    Store: %s\n", getStringValue(ent, "store", "N/A"))
			expiresAt := getStringValue(ent, "expires_at", "")
			if expiresAt != "" && expiresAt != "null" {
				if len(expiresAt) >= 10 {
					fmt.Printf("    Expires: %s\n", expiresAt[:10])
				} else {
					fmt.Printf("    Expires: %s\n", expiresAt)
				}
			} else {
				fmt.Println("    Expires: Never")
			}
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No entitlements found."))
	}

	return nil
}

func runSubscriberSubscriptions(cmd *cobra.Command, args []string) error {
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

	encodedID := url.QueryEscape(appUserID)
	path := fmt.Sprintf("/projects/%s/customers/%s/subscriptions", cfg.ProjectID, encodedID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render(fmt.Sprintf("\n📦 Subscriptions for %s:\n", appUserID)))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			sub := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  Product: %s", sub["product_id"])))
			fmt.Printf("    Store: %s\n", getStringValue(sub, "store", "N/A"))
			fmt.Printf("    Status: %s\n", getStringValue(sub, "status", "N/A"))
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No subscriptions found."))
	}

	return nil
}
