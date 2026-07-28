package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"revenuecat-cli/config"
	rcinternal "revenuecat-cli/internal"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects (Internal Dashboard API)",
	Long: `Manage projects using the internal dashboard API at https://app.revenuecat.com/internal/v1

This is separate from the public /v2 API. The internal API uses session cookie
authentication (rc login) vs the public API's API key authentication.`,
}

var projectsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all projects",
	RunE:    runProjectsList,
}

var projectsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get project details",
	RunE:  runProjectGet,
}

var projectsUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Set default project",
	RunE:  runProjectsUse,
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE:  runProjectsCreate,
}

func init() {
	projectsCmd.AddCommand(projectsListCmd, projectsGetCmd, projectsUseCmd, projectsCreateCmd)

	projectsGetCmd.Flags().StringP("project-id", "i", "", "Project ID")
	projectsUseCmd.Flags().StringP("project-id", "i", "", "Project ID (required)")
	projectsCreateCmd.Flags().StringP("name", "n", "", "Project name (required)")
	projectsCreateCmd.Flags().String("category", "general", "Project category")
	projectsCreateCmd.Flags().StringSlice("platforms", []string{"apple_native"}, "Expected platforms (apple_native, google_play, stripe, amazon, etc.)")
}

func runProjectsList(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	c, err := DashboardNoProject("\n📁 Fetching projects...")
	if err != nil {
		return err
	}

	path := "/developers/me/projects"
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var projects []rcinternal.Project
		if err := json.Unmarshal(ToJSON(resp.Items), &projects); err != nil {
			return fmt.Errorf("error parsing projects: %w", err)
		}

		if len(projects) == 0 {
			fmt.Println(YellowStyle.Render("No projects found."))
			return nil
		}

		fmt.Println(InternalStyle.Render("\n📁 Your Projects:\n"))
		for _, p := range projects {
			currentMarker := ""
			if p.ID == cfg.ProjectID {
				currentMarker = GreenStyle.Render(" (current)")
			}
			fmt.Printf("  ID: %s%s\n", CyanStyle.Render(p.ID), currentMarker)
			fmt.Printf("  Name: %s\n", p.Name)
			fmt.Printf("  Owner: %s\n", p.OwnerEmail)
			if p.RestrictedAccess {
				fmt.Println("  ⚠ Restricted access")
			}
			fmt.Println()
		}
		fmt.Println(GrayStyle.Render(fmt.Sprintf("Total: %d projects", len(projects))))
		fmt.Println(GrayStyle.Render("\nUse 'rc internal projects use -i <project_id>' to set default project"))

		return nil
	})
}
func runProjectGet(cmd *cobra.Command, args []string) error {
	projectID, _ := cmd.Flags().GetString("project-id")
	if projectID == "" {
		return fmt.Errorf("project-id is required (--project-id or -i)")
	}

	c, err := DashboardNoProject("\n📁 Fetching project details...")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/developers/me/projects/%s", projectID)
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var data map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Data), &data); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}

		fmt.Println(InternalStyle.Render("\n📁 Project Details:\n"))
		fmt.Printf("  ID: %s\n", CyanStyle.Render(projectID))
		fmt.Printf("  Name: %s\n", data["name"])
		if owner, ok := data["owner_email"].(string); ok {
			fmt.Printf("  Owner: %s\n", owner)
		}
		if role, ok := data["role"].(string); ok {
			fmt.Printf("  Role: %s\n", role)
		}
		if plan, ok := data["owner_plan"].(string); ok {
			fmt.Printf("  Plan: %s\n", plan)
		}

		return nil
	})
}
func runProjectsUse(cmd *cobra.Command, args []string) error {
	projectID, _ := cmd.Flags().GetString("project-id")
	if projectID == "" {
		return fmt.Errorf("project-id is required (--project-id or -i)")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	cfg.ProjectID = projectID
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	if JSONOutput {
		EmitJSONValue(map[string]string{"project_id": projectID})
		return nil
	}
	fmt.Println(GreenStyle.Render("\n✓ Default project set to: ") + CyanStyle.Render(projectID))
	return nil
}

func runProjectsCreate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}

	category, _ := cmd.Flags().GetString("category")
	platforms, _ := cmd.Flags().GetStringSlice("platforms")

	c, err := DashboardNoProject("\n📁 Creating project...")
	if err != nil {
		return err
	}

	path := "/developers/me/projects"
	data := map[string]interface{}{
		"name":               name,
		"category":           category,
		"expected_platforms": platforms,
	}

	resp, err := c.Client.Post(path, data)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}

	var project map[string]interface{}
	if err := json.Unmarshal(ToJSON(resp.Data), &project); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}
	newID := GetStringValue(project, "id", "")

	// Persist before rendering: this must happen whatever the output mode.
	// Doing it inside the human branch is how --json silently stopped setting
	// the default project.
	setDefault := false
	if newID != "" {
		if cfg, err := config.LoadConfig(); err == nil {
			cfg.ProjectID = newID
			setDefault = config.SaveConfig(cfg) == nil
		}
	}

	return c.Respond(resp, func() error {
		fmt.Println(GreenStyle.Render("\n✓ Project created:"))
		fmt.Printf("  ID: %s\n", CyanStyle.Render(newID))
		fmt.Printf("  Name: %s\n", project["name"])
		if setDefault {
			fmt.Println(GreenStyle.Render("  Set as default project"))
		}
		return nil
	})
}
