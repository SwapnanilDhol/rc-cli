package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

var appsStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
var cyanStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

func initApps() {
	appsCmd := &cobra.Command{
		Use:   "apps",
		Short: "Manage RevenueCat apps",
	}

	appsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all apps",
		RunE:  runListApps,
	}
	appsListCmd.Aliases = []string{"ls"}

	appsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get app details by app_id",
		RunE:  runGetApp,
	}
	appsGetCmd.Flags().StringP("app-id", "i", "", "App ID")

	appsCmd.AddCommand(appsListCmd, appsGetCmd)
	RootCmd.AddCommand(appsCmd)

	// Aliases with colons
	RootCmd.AddCommand(&cobra.Command{
		Use:   "apps:list",
		Short: "List all apps (alias)",
		RunE:  runListApps,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use:   "apps:get",
		Short: "Get app details (alias)",
		RunE:  runGetApp,
	})
}

func runListApps(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n📱 Fetching apps...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/apps", cfg.ProjectID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n📱 Your Apps:\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			app := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", app["id"])))
			fmt.Printf("  Name: %v\n", app["name"])
			fmt.Printf("  Platform: %v\n", app["type"])
			if appStore, ok := app["app_store"].(map[string]interface{}); ok {
				fmt.Printf("  Bundle ID: %v\n", appStore["bundle_id"])
			} else {
				fmt.Println("  Bundle ID: N/A")
			}
			if createdAt, ok := app["created_at"].(string); ok {
				if len(createdAt) >= 10 {
					fmt.Printf("  Created: %s\n", createdAt[:10])
				}
			}
			fmt.Println()
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("Total: %d apps", len(resp.Items))))
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No apps found."))
	}

	return nil
}

func runGetApp(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	appID, _ := cmd.Flags().GetString("app-id")
	if appID == "" {
		return fmt.Errorf("app-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📱 Fetching app details...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s/apps/%s", cfg.ProjectID, appID))
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		fmt.Println(appsStyle.Render("\n📱 App Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", data["id"])))
		fmt.Printf("  Name: %v\n", data["name"])
		fmt.Printf("  Platform: %v\n", data["type"])
		if appStore, ok := data["app_store"].(map[string]interface{}); ok {
			fmt.Printf("  Bundle ID: %v\n", appStore["bundle_id"])
		} else {
			fmt.Println("  Bundle ID: N/A")
		}
		if createdAt, ok := data["created_at"].(string); ok {
			if len(createdAt) >= 10 {
				fmt.Printf("  Created: %s\n", createdAt[:10])
			}
		}
	}

	return nil
}
