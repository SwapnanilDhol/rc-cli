package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initProjects() {
	projectsCmd := &cobra.Command{
		Use:   "projects",
		Short: "Manage RevenueCat projects",
	}

	projectsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		RunE:  runListProjects,
	}
	projectsListCmd.Aliases = []string{"ls"}

	projectsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get project details",
		RunE:  runGetProject,
	}
	projectsGetCmd.Flags().StringP("project-id", "i", "", "Project ID")

	projectsCurrentCmd := &cobra.Command{
		Use:   "current",
		Short: "Get current project details",
		RunE:  runGetCurrentProject,
	}

	projectsCmd.AddCommand(projectsListCmd, projectsGetCmd, projectsCurrentCmd)
	RootCmd.AddCommand(projectsCmd)

	// Aliases
	RootCmd.AddCommand(&cobra.Command{
		Use:   "project:get",
		Short: "Get current project details",
		RunE:  runGetCurrentProject,
	})
}

func runListProjects(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📁 Fetching projects...")

	resp, err := client.Get("/projects")
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n📁 Your Projects:\n"))

	if len(resp.Items) > 0 {
		for _, item := range resp.Items {
			project := item.(map[string]interface{})
			fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", project["id"])))
			fmt.Printf("  Name: %v\n", project["name"])
			if createdAt, ok := project["created_at"].(string); ok {
				if len(createdAt) >= 10 {
					fmt.Printf("  Created: %s\n", createdAt[:10])
				}
			}
			fmt.Println()
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("To set a default project, run: rc config"))
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No projects found."))
	}

	return nil
}

func runGetProject(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	projectID, _ := cmd.Flags().GetString("project-id")
	if projectID == "" {
		return fmt.Errorf("project-id is required")
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📁 Fetching project details...")

	resp, err := client.Get(fmt.Sprintf("/projects/%s", projectID))
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		fmt.Println(appsStyle.Render("\n📁 Project Details:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", data["id"])))
		fmt.Printf("  Name: %v\n", data["name"])
		if createdAt, ok := data["created_at"].(string); ok {
			if len(createdAt) >= 10 {
				fmt.Printf("  Created: %s\n", createdAt[:10])
			}
		}
	}

	return nil
}

func runGetCurrentProject(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	fmt.Println("\n📁 Fetching projects...")

	resp, err := client.Get("/projects")
	if err != nil {
		return err
	}

	if len(resp.Items) == 1 {
		project := resp.Items[0].(map[string]interface{})
		fmt.Println(appsStyle.Render("\n📁 Current Project:\n"))
		fmt.Println(cyanStyle.Render(fmt.Sprintf("  ID: %v", project["id"])))
		fmt.Printf("  Name: %v\n", project["name"])
		if createdAt, ok := project["created_at"].(string); ok {
			if len(createdAt) >= 10 {
				fmt.Printf("  Created: %s\n", createdAt[:10])
			}
		}
	} else if len(resp.Items) > 1 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ Multiple projects found. Specify with:"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  rc config (to set default)"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  Or: rc projects list"))
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("No projects found."))
	}

	return nil
}
