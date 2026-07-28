package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var utilitiesCmd = &cobra.Command{
	Use:   "utilities",
	Short: "Utility endpoints",
}

var utilitiesCountriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "List supported countries",
	RunE:  runUtilitiesCountries,
}

var storesStatusCmd = &cobra.Command{
	Use:     "stores-status",
	Aliases: []string{"stores", "product-stores"},
	Short:   "Get product stores connection status",
	RunE:    runProductStoresStatuses,
}

var collaboratorsCmd = &cobra.Command{
	Use:   "collaborators",
	Short: "Manage collaborators",
}

var collaboratorsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List collaborators",
	RunE:    runCollaboratorsList,
}

var apiKeysCmd = &cobra.Command{
	Use:   "apikeys",
	Short: "Manage API keys",
}

var apiKeysListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List API keys",
	RunE:    runAPIKeysList,
}

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "View audit logs",
}

var auditListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List audit logs",
	RunE:    runAuditList,
}

func init() {
	utilitiesCmd.AddCommand(utilitiesCountriesCmd)

	collaboratorsCmd.AddCommand(collaboratorsListCmd)
	apiKeysCmd.AddCommand(apiKeysListCmd)
	auditCmd.AddCommand(auditListCmd)
}

func runUtilitiesCountries(cmd *cobra.Command, args []string) error {
	c, err := DashboardNoProject(cmd, "\n🌍 Fetching countries...")
	if err != nil {
		return err
	}

	path := "/utilities/countries"
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var countries []map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Items), &countries); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}

		if len(countries) == 0 {
			fmt.Println(YellowStyle.Render("No countries found."))
			return nil
		}

		fmt.Println(InternalStyle.Render("\n🌍 Supported Countries:\n"))
		for _, c := range countries {
			code := c["code"].(string)
			name := c["name"].(string)
			fmt.Printf("  %s - %s\n", CyanStyle.Render(code), name)
		}

		return nil
	})
}
func runProductStoresStatuses(cmd *cobra.Command, args []string) error {
	c, err := Dashboard(cmd, "\n🏪 Fetching product stores statuses...")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/developers/me/projects/%s/product_stores_statuses", c.ProjectID)
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var appStatuses []map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Items), &appStatuses); err != nil {
			return fmt.Errorf("error parsing store statuses response: %w", err)
		}

		if len(appStatuses) == 0 {
			fmt.Println(YellowStyle.Render("No store statuses found."))
			return nil
		}

		fmt.Println(InternalStyle.Render("\n🏪 Product Stores Status:\n"))
		shownAny := false
		for _, appStatus := range appStatuses {
			appID, _ := appStatus["app_config_id"].(string)
			storeStatuses, _ := appStatus["store_statuses"].([]interface{})
			for _, ss := range storeStatuses {
				s := ss.(map[string]interface{})
				productID, _ := s["product_identifier"].(string)
				productTitle, _ := s["product_title"].(string)
				status, _ := s["status"].(string)
				helperText, _ := s["helper_text"].(string)
				storeAppID, _ := s["store_app_id"].(string)
				storeProductID, _ := s["store_product_id"].(string)

				statusColor := GreenStyle
				if status == "needs_action" {
					statusColor = YellowStyle
				} else if status == "in_progress" {
					statusColor = CyanStyle
				}

				fmt.Printf("  App: %s | Product: %s (%s)\n", CyanStyle.Render(appID), productTitle, productID)
				fmt.Printf("    Store App ID: %s | Store Product ID: %s\n", CyanStyle.Render(storeAppID), CyanStyle.Render(storeProductID))
				fmt.Printf("    Status: %s | %s\n", statusColor.Render(status), GrayStyle.Render(helperText))
				fmt.Println()
				shownAny = true
			}
		}

		if !shownAny {
			fmt.Println(YellowStyle.Render("No store statuses found."))
		}

		return nil
	})
}
func runCollaboratorsList(cmd *cobra.Command, args []string) error {
	c, err := Dashboard(cmd, "\n👥 Fetching collaborators...")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/developers/me/projects/%s/collaborators", c.ProjectID)
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var data map[string]interface{}
		var collaborators []interface{}
		var owner map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Data), &data); err != nil {
			// resp.Data is an array directly
			var raw []interface{}
			if err2 := json.Unmarshal(ToJSON(resp.Data), &raw); err2 != nil {
				return fmt.Errorf("error parsing response: %w", err)
			}
			collaborators = raw
		} else {
			collaborators, _ = data["collaborators"].([]interface{})
			owner, _ = data["owner"].(map[string]interface{})
		}

		fmt.Println(InternalStyle.Render("\n👥 Team:\n"))

		if owner != nil {
			fmt.Printf("  Owner: %s\n", CyanStyle.Render(owner["email"].(string)))
		}

		for _, c := range collaborators {
			collab := c.(map[string]interface{})
			fmt.Printf("  %s (%s)\n", collab["email"], collab["role"])
		}

		return nil
	})
}
func runAPIKeysList(cmd *cobra.Command, args []string) error {
	c, err := Dashboard(cmd, "\n🔑 Fetching API keys...")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/developers/me/projects/%s/api_keys", c.ProjectID)
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var keys []map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Items), &keys); err != nil {
			return fmt.Errorf("error parsing keys: %w", err)
		}

		if len(keys) == 0 {
			fmt.Println(YellowStyle.Render("No API keys found."))
			return nil
		}

		fmt.Println(InternalStyle.Render("\n🔑 API Keys:\n"))
		for _, k := range keys {
			keyType := k["key_type"].(string)
			keyStyle := YellowStyle
			if keyType == "secret" {
				keyStyle = CyanStyle
			}
			fmt.Printf("  %s: %s\n", keyStyle.Render(k["label"].(string)), k["id"].(string))
			fmt.Printf("    Type: %s\n", keyType)
			fmt.Printf("    Key: %s\n", k["key"].(string))
			if perms, ok := k["permissions"].([]interface{}); ok && len(perms) > 0 {
				fmt.Printf("    Permissions: %d\n", len(perms))
			}
			fmt.Println()
		}

		return nil
	})
}
func runAuditList(cmd *cobra.Command, args []string) error {
	c, err := Dashboard(cmd, "\n📜 Fetching audit logs...")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/developers/me/projects/%s/audit_logs", c.ProjectID)
	resp, err := c.Client.Get(path)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var data map[string]interface{}
		var logs []interface{}
		if err := json.Unmarshal(ToJSON(resp.Data), &data); err != nil {
			// resp.Data is an array, not a wrapped object
			if err2 := json.Unmarshal(ToJSON(resp.Data), &logs); err2 != nil {
				return fmt.Errorf("error parsing response: %w", err)
			}
		} else {
			logs, _ = data["data"].([]interface{})
		}

		if len(logs) == 0 {
			fmt.Println(YellowStyle.Render("No audit logs found."))
			return nil
		}

		fmt.Println(InternalStyle.Render("\n📜 Recent Audit Logs:\n"))
		for _, l := range logs {
			log := l.(map[string]interface{})
			actionType := log["action_type"].(string)
			targetType := log["target_type"].(string)
			target := log["target_identifier"].(string)
			occurredAt := log["occurred_at"].(string)

			fmt.Printf("  %s %s %s %s\n",
				CyanStyle.Render(actionType),
				GrayStyle.Render(targetType),
				YellowStyle.Render(target),
				GrayStyle.Render(occurredAt[:10]))
		}

		return nil
	})
}
