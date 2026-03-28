package cmd

import (
	"fmt"
	"net/url"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initEntitlements() {
	entitlementsCmd := &cobra.Command{
		Use:   "entitlements",
		Short: "Manage entitlements",
	}

	entitlementsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all entitlements",
		RunE:  runListEntitlements,
	}

	entitlementsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get entitlement details",
		RunE:  runGetEntitlement,
	}
	entitlementsGetCmd.Flags().StringP("entitlement-id", "i", "", "Entitlement ID")

	entitlementsProductsCmd := &cobra.Command{
		Use:   "products",
		Short: "Get products in an entitlement",
		RunE:  runEntitlementProducts,
	}
	entitlementsProductsCmd.Flags().StringP("entitlement-id", "i", "", "Entitlement ID")

	entitlementsActiveCmd := &cobra.Command{
		Use:   "active",
		Short: "Get active entitlements for a customer",
		RunE:  runActiveEntitlements,
	}
	entitlementsActiveCmd.Flags().StringP("app-user-id", "u", "", "App User ID")

	entitlementsCmd.AddCommand(entitlementsListCmd, entitlementsGetCmd, entitlementsProductsCmd, entitlementsActiveCmd)
	RootCmd.AddCommand(entitlementsCmd)

	// Alias
	RootCmd.AddCommand(&cobra.Command{
		Use:   "entitlements:list",
		Short: "List entitlements",
		RunE:  runListEntitlements,
	})
}

func runListEntitlements(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n📋 Fetching entitlements...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/entitlements", cfg.ProjectID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n📋 Entitlements:\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			ent := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", ent["id"])))
			fmt.Printf("    Name: %s\n", ent["display_name"])
			fmt.Printf("    Lookup Key: %s\n", ent["lookup_key"])
			fmt.Printf("    State: %s\n", ent["state"])
			if createdAt, ok := ent["created_at"].(float64); ok {
				fmt.Printf("    Created: %s\n", formatUnixTime(int64(createdAt)))
			}
			fmt.Println()
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("Total: %d entitlements", len(resp.Items))))
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No entitlements found."))
	}

	return nil
}

func runGetEntitlement(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Fetching entitlement details...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/entitlements/%s", cfg.ProjectID, entitlementID))
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		fmt.Println(appsStyle.Render("\n📋 Entitlement Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", data["id"])))
		fmt.Printf("  Name: %s\n", data["display_name"])
		fmt.Printf("  Lookup Key: %s\n", data["lookup_key"])
		fmt.Printf("  State: %s\n", data["state"])
		if createdAt, ok := data["created_at"].(float64); ok {
			fmt.Printf("  Created: %s\n", formatUnixTime(int64(createdAt)))
		}

		// Get products
		productsResp, err := client.Get(fmt.Sprintf("/projects/%s/entitlements/%s/products", cfg.ProjectID, entitlementID))
		if err == nil && len(productsResp.Items) > 0 {
			fmt.Println(appsStyle.Render("\n  Products:"))
			for _, p := range productsResp.Items {
				product := p.(map[string]interface{})
				fmt.Printf("    - %s (%s)\n", product["id"], product["store"])
			}
		}
	}

	return nil
}

func runEntitlementProducts(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n🛍️ Fetching products...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/entitlements/%s/products", cfg.ProjectID, entitlementID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render(fmt.Sprintf("\n🛍️ Products for Entitlement %s:\n", entitlementID)))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			product := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", product["id"])))
			fmt.Printf("    Store: %s\n", product["store"])
			fmt.Printf("    Type: %s\n", product["type"])
			fmt.Printf("    Store ID: %s\n", getStringValue(product, "store_identifier", "N/A"))

			if sub, ok := product["subscription"].(map[string]interface{}); ok {
				if duration, ok := sub["duration"].(string); ok && duration != "" {
					fmt.Printf("    Duration: %s\n", duration)
				}
			}
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No products found."))
	}

	return nil
}

func runActiveEntitlements(cmd *cobra.Command, args []string) error {
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

	// Use the customer endpoint directly for active entitlements
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
