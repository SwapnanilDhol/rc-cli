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
		Short: "Promotional offers (public v2 API)",
	}

	offersListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all promotional offers",
		RunE:    runListOffers,
	}

	offersCmd.AddCommand(offersListCmd)
	RootCmd.AddCommand(offersCmd)
}

func runListOffers(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		return fmt.Errorf("no project ID configured. Run: rc config, or pass --project-id")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	resp, err := client.Get(fmt.Sprintf("/projects/%s/offers", cfg.ProjectID))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s %s", resp.StatusCode, resp.Error, resp.Message)
	}

	if jsonOutput {
		return emitJSON(resp.Items)
	}

	fmt.Println(appsStyle.Render("\n🏷️ Promotional Offers\n"))
	if len(resp.Items) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  No promotional offers found.\n"))
		return nil
	}
	for _, item := range resp.Items {
		offer, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %s", offer["id"])))
		fmt.Printf("  Name: %s\n", offer["name"])
		fmt.Printf("  Type: %s\n", offer["type"])
		fmt.Println()
	}

	return nil
}
