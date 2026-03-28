package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initProducts() {
	productsCmd := &cobra.Command{
		Use:   "products",
		Short: "Manage products",
	}

	productsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all products",
		RunE:  runListProducts,
	}
	productsListCmd.Aliases = []string{"ls"}

	productsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get product details",
		RunE:  runGetProduct,
	}
	productsGetCmd.Flags().StringP("product-id", "i", "", "Product ID")

	productsCmd.AddCommand(productsListCmd, productsGetCmd)
	RootCmd.AddCommand(productsCmd)

	// Offerings subcommand
	offeringsCmd := &cobra.Command{
		Use:   "offerings",
		Short: "Manage offerings",
	}

	offeringsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all offerings",
		RunE:  runListOfferings,
	}
	offeringsListCmd.Aliases = []string{"ls"}

	offeringsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get offering details",
		RunE:  runGetOffering,
	}
	offeringsGetCmd.Flags().StringP("offering-id", "i", "", "Offering ID")

	offeringsPackagesCmd := &cobra.Command{
		Use:   "packages",
		Short: "List packages in an offering",
		RunE:  runOfferingPackages,
	}
	offeringsPackagesCmd.Flags().StringP("offering-id", "i", "", "Offering ID")

	offeringsCmd.AddCommand(offeringsListCmd, offeringsGetCmd, offeringsPackagesCmd)
	RootCmd.AddCommand(offeringsCmd)

	// Packages subcommand
	packagesCmd := &cobra.Command{
		Use:   "packages",
		Short: "Manage packages",
	}

	packagesGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get package details",
		RunE:  runGetPackage,
	}
	packagesGetCmd.Flags().StringP("package-id", "i", "", "Package ID")

	packagesProductsCmd := &cobra.Command{
		Use:   "products",
		Short: "Get products in a package",
		RunE:  runPackageProducts,
	}
	packagesProductsCmd.Flags().StringP("package-id", "i", "", "Package ID")

	packagesCmd.AddCommand(packagesGetCmd, packagesProductsCmd)
	RootCmd.AddCommand(packagesCmd)

	// Aliases
	RootCmd.AddCommand(&cobra.Command{
		Use:   "packages:list",
		Short: "List offerings (alias)",
		RunE:  runListOfferings,
	})
}

func runListProducts(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n🛍️ Fetching products...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/products", cfg.ProjectID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n🛍️ Products:\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			product := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", product["id"])))
			fmt.Printf("    Store: %s\n", product["store"])
			fmt.Printf("    Store ID: %s\n", getStringValue(product, "store_identifier", "N/A"))
			fmt.Printf("    Type: %s\n", getStringValue(product, "type", "N/A"))

			if sub, ok := product["subscription"].(map[string]interface{}); ok {
				if duration, ok := sub["duration"].(string); ok && duration != "" {
					fmt.Printf("    Duration: %s\n", duration)
				}
				if trial, ok := sub["trial_duration"].(string); ok && trial != "" {
					fmt.Printf("    Trial: %s\n", trial)
				}
			}
			fmt.Println()
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("Total: %d products", len(resp.Items))))
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No products found."))
	}

	return nil
}

func runGetProduct(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	productID, _ := cmd.Flags().GetString("product-id")
	if productID == "" {
		return fmt.Errorf("product-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n🛍️ Fetching product details...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/products/%s", cfg.ProjectID, productID))
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		fmt.Println(appsStyle.Render("\n🛍️ Product Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", data["id"])))
		fmt.Printf("  Store: %s\n", getStringValue(data, "store", "N/A"))
		fmt.Printf("  Store ID: %s\n", getStringValue(data, "store_identifier", "N/A"))
		fmt.Printf("  Type: %s\n", getStringValue(data, "type", "N/A"))
		fmt.Printf("  State: %s\n", getStringValue(data, "state", "N/A"))

		if sub, ok := data["subscription"].(map[string]interface{}); ok {
			fmt.Println(appsStyle.Render("\n  Subscription:"))
			if duration, ok := sub["duration"].(string); ok && duration != "" {
				fmt.Printf("    Duration: %s\n", duration)
			}
			if trial, ok := sub["trial_duration"].(string); ok && trial != "" {
				fmt.Printf("    Trial Duration: %s\n", trial)
			}
			if grace, ok := sub["grace_period_duration"].(string); ok && grace != "" {
				fmt.Printf("    Grace Period: %s\n", grace)
			}
		}
	}

	return nil
}

func runListOfferings(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n📦 Fetching offerings...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/offerings", cfg.ProjectID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n📦 Offerings:\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			offering := item.(map[string]interface{})
			isCurrent := getBoolValue(offering, "is_current")
			status := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Current")
			if !isCurrent {
				status = "Inactive"
			}

			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", offering["id"])))
			fmt.Printf("  Name: %v\n", offering["display_name"])
			fmt.Printf("  Lookup Key: %v\n", offering["lookup_key"])
			fmt.Printf("  Status: %s\n", status)
			fmt.Printf("  State: %v\n", offering["state"])
			if createdAt, ok := offering["created_at"].(float64); ok {
				fmt.Printf("  Created: %s\n", formatUnixTime(int64(createdAt)))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No offerings found."))
	}

	return nil
}

func runGetOffering(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	offeringID, _ := cmd.Flags().GetString("offering-id")
	if offeringID == "" {
		return fmt.Errorf("offering-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Fetching offering details...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/offerings/%s", cfg.ProjectID, offeringID))
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		isCurrent := getBoolValue(data, "is_current")
		status := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Current")
		if !isCurrent {
			status = "Inactive"
		}

		fmt.Println(appsStyle.Render("\n📦 Offering Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", data["id"])))
		fmt.Printf("  Name: %v\n", data["display_name"])
		fmt.Printf("  Lookup Key: %v\n", data["lookup_key"])
		fmt.Printf("  Status: %s\n", status)
		fmt.Printf("  State: %v\n", data["state"])
		if createdAt, ok := data["created_at"].(float64); ok {
			fmt.Printf("  Created: %s\n", formatUnixTime(int64(createdAt)))
		}

		// Get packages
		pkgsResp, err := client.Get(fmt.Sprintf("/projects/%s/offerings/%s/packages", cfg.ProjectID, offeringID))
		if err == nil && len(pkgsResp.Items) > 0 {
			fmt.Println(appsStyle.Render("\n  Packages:"))
			for _, p := range pkgsResp.Items {
				pkg := p.(map[string]interface{})
				fmt.Printf("    - %s (%s)\n", pkg["display_name"], pkg["id"])
			}
		}
	}

	return nil
}

func runOfferingPackages(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	offeringID, _ := cmd.Flags().GetString("offering-id")
	if offeringID == "" {
		return fmt.Errorf("offering-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Fetching packages...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/offerings/%s/packages", cfg.ProjectID, offeringID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render(fmt.Sprintf("\n📦 Packages in %s:\n", offeringID)))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			pkg := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", pkg["id"])))
			fmt.Printf("    Name: %s\n", pkg["display_name"])
			fmt.Printf("    Lookup Key: %s\n", pkg["lookup_key"])
			if pos, ok := pkg["position"].(float64); ok {
				fmt.Printf("    Position: %.0f\n", pos)
			}
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No packages found."))
	}

	return nil
}

func runGetPackage(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	packageID, _ := cmd.Flags().GetString("package-id")
	if packageID == "" {
		return fmt.Errorf("package-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Fetching package details...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/packages/%s", cfg.ProjectID, packageID))
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		fmt.Println(appsStyle.Render("\n📋 Package Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", data["id"])))
		fmt.Printf("  Name: %s\n", data["display_name"])
		fmt.Printf("  Lookup Key: %s\n", data["lookup_key"])
		if pos, ok := data["position"].(float64); ok {
			fmt.Printf("  Position: %.0f\n", pos)
		}
		if createdAt, ok := data["created_at"].(float64); ok {
			fmt.Printf("  Created: %s\n", formatUnixTime(int64(createdAt)))
		}
	}

	return nil
}

func runPackageProducts(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	packageID, _ := cmd.Flags().GetString("package-id")
	if packageID == "" {
		return fmt.Errorf("package-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n🛍️ Fetching products...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/packages/%s/products", cfg.ProjectID, packageID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render(fmt.Sprintf("\n🛍️ Products in %s:\n", packageID)))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			if pkgProduct, ok := item.(map[string]interface{}); ok {
				if eligibility, ok := pkgProduct["eligibility_criteria"].(string); ok {
					fmt.Printf("  Eligibility: %s\n", eligibility)
				}
				if product, ok := pkgProduct["product"].(map[string]interface{}); ok {
					fmt.Println(cyanStyle.Render(fmt.Sprintf("  Product: %s", product["id"])))
					fmt.Printf("    Store: %s\n", product["store"])
					fmt.Printf("    Type: %s\n", product["type"])

					if sub, ok := product["subscription"].(map[string]interface{}); ok {
						if duration, ok := sub["duration"].(string); ok && duration != "" {
							fmt.Printf("    Duration: %s\n", duration)
						}
					}
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No products found."))
	}

	return nil
}
