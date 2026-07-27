package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	rcinternal "revenuecat-cli/internal"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps (Internal Dashboard API)",
	Long:  `Manage apps using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var appsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all apps",
	RunE:    runAppsList,
}

var appsSubscriptionGroupsCmd = &cobra.Command{
	Use:   "subscription-groups",
	Short: "List App Store Connect subscription groups for an app",
	RunE:  runAppsSubscriptionGroups,
}

var appStoreProductsCmd = &cobra.Command{
	Use:   "app-store-products",
	Short: "Manage App Store Connect products via app_store_products",
}

var appStoreProductsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an app store product",
	RunE:  runAppStoreProductsCreate,
}

func init() {
	appsCmd.AddCommand(appsListCmd, appsSubscriptionGroupsCmd, appStoreProductsCmd)
	appStoreProductsCmd.AddCommand(appStoreProductsCreateCmd)

	appsSubscriptionGroupsCmd.Flags().StringP("app-id", "i", "", "App ID (required; e.g. appb5d7b1196a)")

	appStoreProductsCreateCmd.Flags().StringP("app-id", "i", "", "App ID (required; e.g. appb5d7b1196a)")
	appStoreProductsCreateCmd.Flags().String("product-type", "", "Product type (e.g. subscriptions or inAppPurchases)")
	appStoreProductsCreateCmd.Flags().String("in-app-purchase-type", "", "In-app purchase type (e.g. NON_RENEWING_SUBSCRIPTION)")
	appStoreProductsCreateCmd.Flags().String("duration", "", "Subscription duration (ONE_WEEK|ONE_MONTH|TWO_MONTHS|THREE_MONTHS|SIX_MONTHS|ONE_YEAR)")
	appStoreProductsCreateCmd.Flags().String("subscription-group-id", "", "App Store Connect subscription group id (required for subscriptions + duration)")
	appStoreProductsCreateCmd.Flags().String("subscription-group-name", "", "App Store Connect subscription group name (optional)")
	appStoreProductsCreateCmd.Flags().String("identifier", "", "Product identifier / product_identifier (required)")
	appStoreProductsCreateCmd.Flags().StringP("name", "n", "", "Display name / name (required)")
}

func runAppsList(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📱 Fetching apps...")

	path := fmt.Sprintf("/developers/me/projects/%s/apps", projectID)
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

	var apps []rcinternal.App
	if err := json.Unmarshal(ToJSON(resp.Items), &apps); err != nil {
		return fmt.Errorf("error parsing apps: %w", err)
	}

	if len(apps) == 0 {
		fmt.Println(YellowStyle.Render("No apps found."))
		return nil
	}

	fmt.Println(InternalStyle.Render("\n📱 Apps:\n"))
	for _, a := range apps {
		fmt.Printf("  ID: %s\n", CyanStyle.Render(a.ID))
		fmt.Printf("  Name: %s\n", a.Name)
		fmt.Printf("  Type: %s\n", a.Type)
		fmt.Println()
	}

	return nil
}

func runAppsSubscriptionGroups(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	appID, _ := cmd.Flags().GetString("app-id")
	if appID == "" {
		return fmt.Errorf("--app-id is required")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📚 Fetching subscription groups for app...")

	path := fmt.Sprintf("/developers/me/projects/%s/apps/%s/subscription_groups", projectID, appID)
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

	var groups []map[string]interface{}
	if len(resp.Items) > 0 {
		_ = json.Unmarshal(ToJSON(resp.Items), &groups)
	}
	if len(groups) == 0 && resp.Data != nil {
		var asArray []map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Data), &asArray); err == nil {
			groups = asArray
		} else {
			var asOne map[string]interface{}
			if err2 := json.Unmarshal(ToJSON(resp.Data), &asOne); err2 == nil {
				groups = []map[string]interface{}{asOne}
			}
		}
	}

	if len(groups) == 0 {
		fmt.Println(YellowStyle.Render("No subscription groups found."))
		return nil
	}

	fmt.Println(InternalStyle.Render("\n📚 Subscription Groups:\n"))
	for _, g := range groups {
		id := GetStringValue(g, "id", GetStringValue(g, "subscription_group_id", ""))
		identifier := GetStringValue(g, "identifier", GetStringValue(g, "subscription_group_identifier", ""))
		name := GetStringValue(g, "name", GetStringValue(g, "display_name", GetStringValue(g, "group_name", "")))

		fmt.Printf("  ID: %s\n", CyanStyle.Render(id))
		if identifier != "" {
			fmt.Printf("  Identifier: %s\n", identifier)
		}
		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		if products, ok := g["products"].([]interface{}); ok {
			fmt.Printf("  Products: %d\n", len(products))
		}
		fmt.Println()
	}

	return nil
}

func validateDurationEnum(duration string) error {
	switch duration {
	case "ONE_WEEK", "ONE_MONTH", "TWO_MONTHS", "THREE_MONTHS", "SIX_MONTHS", "ONE_YEAR":
		return nil
	default:
		return fmt.Errorf("--duration must be one of: ONE_WEEK, ONE_MONTH, TWO_MONTHS, THREE_MONTHS, SIX_MONTHS, ONE_YEAR (got %q)", duration)
	}
}

func runAppStoreProductsCreate(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	appID, _ := cmd.Flags().GetString("app-id")
	if appID == "" {
		return fmt.Errorf("--app-id is required")
	}

	productType, _ := cmd.Flags().GetString("product-type")
	if productType == "" {
		return fmt.Errorf("--product-type is required")
	}

	inAppPurchaseType, _ := cmd.Flags().GetString("in-app-purchase-type")
	duration, _ := cmd.Flags().GetString("duration")
	subscriptionGroupID, _ := cmd.Flags().GetString("subscription-group-id")
	subscriptionGroupName, _ := cmd.Flags().GetString("subscription-group-name")
	identifier, _ := cmd.Flags().GetString("identifier")
	displayName, _ := cmd.Flags().GetString("name")

	if identifier == "" {
		return fmt.Errorf("--identifier is required")
	}
	if displayName == "" {
		return fmt.Errorf("--name is required")
	}

	if duration != "" {
		if err := validateDurationEnum(duration); err != nil {
			return err
		}
		if subscriptionGroupID == "" {
			return fmt.Errorf("--subscription-group-id is required when --duration is set")
		}
	}

	product := map[string]interface{}{
		"product_identifier": identifier,
		"name":               displayName,
		"product_type":       productType,
	}
	if inAppPurchaseType != "" {
		product["in_app_purchase_type"] = inAppPurchaseType
	}
	if duration != "" {
		product["duration"] = duration
	}
	if subscriptionGroupID != "" {
		subGroup := map[string]interface{}{"id": subscriptionGroupID}
		if subscriptionGroupName != "" {
			subGroup["name"] = subscriptionGroupName
		}
		product["subscription_group"] = subGroup
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n📦 Creating app store product (internal POST)...")
	path := fmt.Sprintf("/developers/me/projects/%s/apps/%s/app_store_products", projectID, appID)
	body := map[string]interface{}{"products": []interface{}{product}}

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

	fmt.Println(GreenStyle.Render("\n✓ App store product created."))
	return nil
}
