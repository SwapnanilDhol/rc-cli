package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initOffers() {
	offersCmd := &cobra.Command{
		Use:   "offers",
		Short: "Manage promotional offers",
	}

	offersListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all promotional offers",
		RunE:  runListOffers,
	}

	offersCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new promotional offer",
		RunE:  runCreateOffer,
	}

	offersCmd.AddCommand(offersListCmd, offersCreateCmd)
	RootCmd.AddCommand(offersCmd)

	// Promotions subcommand
	promotionsCmd := &cobra.Command{
		Use:   "promotions",
		Short: "Manage promotions",
	}

	promotionsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all promotions",
		RunE:  runListPromotions,
	}

	promotionsCmd.AddCommand(promotionsListCmd)
	RootCmd.AddCommand(promotionsCmd)

	// Aliases
	RootCmd.AddCommand(&cobra.Command{
		Use:   "intro-offers:list",
		Short: "List introductory offers",
		RunE:  runListOffers,
	})
}

func runListOffers(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n🏷️ Fetching promotional offers...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/offers", cfg.ProjectID))
	if err != nil {
		fmt.Println(appsStyle.Render("\n🏷️ Promotional Offers\n"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Promotional offers are managed in the RevenueCat dashboard"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/offers\n"))
		return nil
	}

	fmt.Println(appsStyle.Render("\n🏷️ Promotional Offers\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			offer := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", offer["id"])))
			fmt.Printf("  Name: %s\n", offer["name"])
			fmt.Printf("  Type: %s\n", offer["type"])
			fmt.Println()
		}
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  No promotional offers found."))
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/offers to create one\n"))
	}

	return nil
}

func runCreateOffer(cmd *cobra.Command, args []string) error {
	fmt.Println(appsStyle.Render("\n🏷️ Create Promotional Offer\n"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Creating promotional offers is done in the RevenueCat dashboard"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/offers\n"))

	return nil
}

func runListPromotions(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n🎁 Fetching promotions...")

	_, err = client.Get(fmt.Sprintf("/projects/%s/promotions", cfg.ProjectID))
	if err != nil {
		fmt.Println(appsStyle.Render("\n🎁 Promotions\n"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Promotions are managed in the RevenueCat dashboard"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Visit: https://app.revenuecat.com/promotions\n"))
		return nil
	}

	return nil
}
