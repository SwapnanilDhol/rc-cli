package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage subscriber lists (Internal Dashboard API)",
	Long:  `Manage subscriber lists using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var listsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List subscriber lists",
	RunE:    runListsList,
}

var listsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get subscriber list details",
	RunE:  runListsGet,
}

var listsManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Get manifest of all subscriber lists",
	RunE:  runListsManifest,
}

func init() {
	listsCmd.AddCommand(listsListCmd, listsGetCmd, listsManifestCmd)

	listsListCmd.Flags().Int("limit", 100, "Maximum number of lists to fetch")
	listsGetCmd.Flags().StringP("list-id", "l", "", "List ID (required)")
}

func runListsList(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")

	Progress("\n👥 Fetching subscriber lists...")

	path := fmt.Sprintf("/developers/me/projects/%s/subscriber_lists", projectID)
	params := map[string]string{"limit": fmt.Sprintf("%d", limit)}
	resp, err := client.GetWithParams(path, params)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(ToJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	lists, _ := data["data"].([]interface{})

	if len(lists) == 0 {
		fmt.Println(YellowStyle.Render("No subscriber lists found."))
		return nil
	}

	fmt.Println(InternalStyle.Render("\n👥 Subscriber Lists:\n"))
	for _, l := range lists {
		list := l.(map[string]interface{})
		id := list["id"].(string)
		name := list["name"].(string)
		size, _ := list["size"].(float64)
		listType, _ := list["type"].(string)

		fmt.Printf("  ID: %s\n", CyanStyle.Render(id))
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Type: %s\n", GreenStyle.Render(listType))
		fmt.Printf("  Size: %.0f\n", size)
		fmt.Println()
	}

	return nil
}

func runListsGet(cmd *cobra.Command, args []string) error {
	listID, _ := cmd.Flags().GetString("list-id")
	if listID == "" {
		return fmt.Errorf("list-id is required (--list-id or -l)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n👥 Fetching subscriber list...")

	path := fmt.Sprintf("/developers/me/subscriber_lists/%s", listID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	var list map[string]interface{}
	if err := json.Unmarshal(ToJSON(resp.Data), &list); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(InternalStyle.Render("\n👥 Subscriber List Details:\n"))
	fmt.Printf("  ID: %s\n", CyanStyle.Render(list["id"].(string)))
	fmt.Printf("  Name: %s\n", list["name"])
	fmt.Printf("  Type: %s\n", GreenStyle.Render(list["type"].(string)))
	if size, ok := list["size"].(float64); ok {
		fmt.Printf("  Size: %.0f\n", size)
	}

	return nil
}

func runListsManifest(cmd *cobra.Command, args []string) error {
	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n👥 Fetching subscriber lists manifest...")

	path := "/developers/me/subscriber_lists/manifest"
	resp, err := client.Get(path)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(ToJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(InternalStyle.Render("\n👥 Subscriber Lists Manifest:\n"))

	if totals, ok := data["totals"].(map[string]interface{}); ok {
		fmt.Println("  Totals:")
		for key, value := range totals {
			fmt.Printf("    %s: %v\n", CyanStyle.Render(key), value)
		}
	}

	if lists, ok := data["lists"].([]interface{}); ok && len(lists) > 0 {
		fmt.Println(InternalStyle.Render("\n  Lists:"))
		for _, l := range lists {
			list := l.(map[string]interface{})
			fmt.Printf("    - %s (%s): %.0f\n",
				CyanStyle.Render(list["id"].(string)),
				list["name"],
				list["size"].(float64))
		}
	}

	return nil
}
