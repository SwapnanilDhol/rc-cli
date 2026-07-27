package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	rcinternal "revenuecat-cli/internal"
)

var entitlementsCmd = &cobra.Command{
	Use:   "entitlements",
	Short: "Manage entitlements (Internal Dashboard API)",
	Long:  `Manage entitlements using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var entitlementsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all entitlements",
	RunE:    runEntitlementsList,
}

var entitlementsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an entitlement",
	RunE:  runEntitlementsCreate,
}

var entitlementsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an entitlement",
	RunE:  runEntitlementsDelete,
}

var entitlementsArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive an entitlement",
	RunE:  runEntitlementsArchive,
}

var entitlementsAttachCmd = &cobra.Command{
	Use:   "attach-products",
	Short: "Attach products to an entitlement",
	RunE:  runEntitlementsAttachProducts,
}

var entitlementsDetachCmd = &cobra.Command{
	Use:   "detach-products",
	Short: "Detach products from an entitlement",
	RunE:  runEntitlementsDetachProducts,
}

func init() {
	entitlementsCmd.AddCommand(entitlementsListCmd, entitlementsCreateCmd, entitlementsDeleteCmd, entitlementsArchiveCmd, entitlementsAttachCmd, entitlementsDetachCmd)

	entitlementsCreateCmd.Flags().StringP("identifier", "i", "", "Entitlement identifier (slug)")
	entitlementsCreateCmd.Flags().StringP("name", "n", "", "Display name")

	entitlementsDeleteCmd.Flags().StringP("entitlement-id", "e", "", "Entitlement ID")
	entitlementsArchiveCmd.Flags().StringP("entitlement-id", "e", "", "Entitlement ID")
	entitlementsAttachCmd.Flags().StringP("entitlement-id", "e", "", "Entitlement ID (required)")
	entitlementsAttachCmd.Flags().StringSlice("product-ids", []string{}, "Product IDs to attach (required)")
	entitlementsDetachCmd.Flags().StringP("entitlement-id", "e", "", "Entitlement ID (required)")
	entitlementsDetachCmd.Flags().StringSlice("product-ids", []string{}, "Product IDs to detach (required)")
}

func runEntitlementsList(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📋 Fetching entitlements...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements", projectID)
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

	var entitlements []rcinternal.Entitlement
	if err := json.Unmarshal(ToJSON(resp.Items), &entitlements); err != nil {
		return fmt.Errorf("error parsing entitlements: %w", err)
	}

	if len(entitlements) == 0 {
		fmt.Println(YellowStyle.Render("No entitlements found."))
		return nil
	}

	fmt.Println(InternalStyle.Render("\n📋 Entitlements:\n"))
	for _, e := range entitlements {
		fmt.Printf("  ID: %s\n", CyanStyle.Render(e.ID))
		fmt.Printf("  Identifier: %s\n", e.Identifier)
		fmt.Printf("  Name: %s\n", e.DisplayName)
		fmt.Printf("  Products: %d\n", len(e.Products))
		fmt.Println()
	}
	fmt.Println(GrayStyle.Render(fmt.Sprintf("Total: %d entitlements", len(entitlements))))

	return nil
}

func runEntitlementsCreate(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	identifier, _ := cmd.Flags().GetString("identifier")
	name, _ := cmd.Flags().GetString("name")

	if identifier == "" {
		return fmt.Errorf("identifier is required (--identifier or -i)")
	}
	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📋 Creating entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements", projectID)
	data := map[string]string{
		"identifier":   identifier,
		"display_name": name,
	}
	resp, err := client.Post(path, data)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	var entitlement rcinternal.Entitlement
	if err := json.Unmarshal(ToJSON(resp.Data), &entitlement); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(GreenStyle.Render("\n✓ Entitlement created:"))
	fmt.Printf("  ID: %s\n", CyanStyle.Render(entitlement.ID))
	fmt.Printf("  Identifier: %s\n", entitlement.Identifier)
	fmt.Printf("  Name: %s\n", entitlement.DisplayName)

	return nil
}

func runEntitlementsDelete(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required (--entitlement-id or -e)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📋 Deleting entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements/%s", projectID, entitlementID)
	resp, err := client.Delete(path)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Entitlement deleted: " + entitlementID))
	return nil
}

func runEntitlementsArchive(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required (--entitlement-id or -e)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📋 Archiving entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements/%s/actions/archive", projectID, entitlementID)
	resp, err := client.Post(path, map[string]interface{}{})
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Entitlement archived"))

	return nil
}

func runEntitlementsAttachProducts(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required (--entitlement-id or -e)")
	}

	productIDs, _ := cmd.Flags().GetStringSlice("product-ids")
	if len(productIDs) == 0 {
		return fmt.Errorf("--product-ids is required")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📎 Attaching products to entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements/%s/attach_products", projectID, entitlementID)
	body := map[string]interface{}{
		"products_ids": productIDs,
	}

	resp, err := client.Post(path, body)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Products attached"))
	return nil
}

func runEntitlementsDetachProducts(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required (--entitlement-id or -e)")
	}

	productIDs, _ := cmd.Flags().GetStringSlice("product-ids")
	if len(productIDs) == 0 {
		return fmt.Errorf("--product-ids is required")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📎 Detaching products from entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements/%s/detach_products", projectID, entitlementID)
	body := map[string]interface{}{
		"products_ids": productIDs,
	}

	resp, err := client.Post(path, body)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Products detached"))
	return nil
}
