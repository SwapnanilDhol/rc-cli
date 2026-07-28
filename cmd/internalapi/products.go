package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	rcinternal "revenuecat-cli/internal"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Manage products (Internal Dashboard API)",
	Long:  `Manage products using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var productsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all products",
	RunE:    runProductsList,
}

var productsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a product for an app",
	RunE:  runProductsCreate,
}

var productsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a product by id",
	RunE:  runProductsUpdate,
}

func init() {
	productsCmd.AddCommand(productsListCmd, productsCreateCmd, productsUpdateCmd)

	productsListCmd.Flags().Int("limit", 100, "Maximum number of products to fetch")
	productsCreateCmd.Flags().String("app-id", "", "App ID (required; e.g. appb5d7b1196a)")
	productsCreateCmd.Flags().String("product-type", "", "Product type (subscription|non_consumable_product|consumable_product|non_renewing_subscription)")
	productsCreateCmd.Flags().StringP("identifier", "i", "", "Product identifier (required)")
	productsCreateCmd.Flags().StringP("name", "n", "", "Display name (display_name) (required)")
	productsUpdateCmd.Flags().String("product-id", "", "Product ID (required)")
	productsUpdateCmd.Flags().String("product-type", "", "Product type (subscription|non_consumable_product|consumable_product|non_renewing_subscription)")
	productsUpdateCmd.Flags().StringP("identifier", "i", "", "Product identifier")
	productsUpdateCmd.Flags().StringP("name", "n", "", "Display name (display_name)")
}

func validateProductType(productType string) error {
	switch productType {
	case "subscription", "non_consumable_product", "consumable_product", "non_renewing_subscription":
		return nil
	default:
		return fmt.Errorf("--product-type must be one of: subscription, non_consumable_product, consumable_product, non_renewing_subscription (got %q)", productType)
	}
}

func runProductsList(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")

	c, err := Dashboard("\n🛍️ Fetching products...")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/developers/me/projects/%s/products", c.ProjectID)
	params := map[string]string{"limit": fmt.Sprintf("%d", limit)}
	resp, err := c.Client.GetWithParams(path, params)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var products []rcinternal.Product

		// Try resp.Items first (some endpoints return products in items array)
		if err := json.Unmarshal(ToJSON(resp.Items), &products); err != nil || len(products) == 0 {
			// Fall back to resp.Data which may have products at data["products"] or data["data"]
			var data map[string]interface{}
			if err2 := json.Unmarshal(ToJSON(resp.Data), &data); err2 == nil {
				if prods, ok := data["products"].([]interface{}); ok {
					itemsJSON, _ := json.Marshal(prods)
					json.Unmarshal(itemsJSON, &products)
				} else if prods, ok := data["data"].([]interface{}); ok {
					itemsJSON, _ := json.Marshal(prods)
					json.Unmarshal(itemsJSON, &products)
				}
			}
		}

		if len(products) == 0 {
			fmt.Println(YellowStyle.Render("No products found."))
			return nil
		}

		fmt.Println(InternalStyle.Render("\n🛍️ Products:\n"))
		for _, p := range products {
			fmt.Printf("  ID: %s\n", CyanStyle.Render(p.ID))
			fmt.Printf("  Identifier: %s\n", p.Identifier)
			fmt.Printf("  Type: %s\n", p.ProductType)
			if p.App != nil {
				fmt.Printf("  App: %s\n", p.App.Name)
			}
			fmt.Println()
		}
		fmt.Println(GrayStyle.Render(fmt.Sprintf("Total: %d products", len(products))))

		return nil
	})
}
func runProductsCreate(cmd *cobra.Command, args []string) error {
	appID, _ := cmd.Flags().GetString("app-id")
	productType, _ := cmd.Flags().GetString("product-type")
	identifier, _ := cmd.Flags().GetString("identifier")
	displayName, _ := cmd.Flags().GetString("name")

	if appID == "" {
		return fmt.Errorf("--app-id is required")
	}
	if productType == "" {
		return fmt.Errorf("--product-type is required")
	}
	if err := validateProductType(productType); err != nil {
		return err
	}
	if identifier == "" {
		return fmt.Errorf("--identifier is required")
	}
	if displayName == "" {
		return fmt.Errorf("--name is required")
	}

	c, err := Dashboard("\n🛍️ Creating product...")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/developers/me/projects/%s/apps/%s/products", c.ProjectID, appID)
	data := map[string]interface{}{
		"product_type": productType,
		"identifier":   identifier,
		"display_name": displayName,
	}

	resp, err := c.Client.Post(path, data)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		var product rcinternal.Product
		if resp.Data != nil {
			_ = json.Unmarshal(ToJSON(resp.Data), &product)
		}
		if product.ID != "" {
			fmt.Println(GreenStyle.Render("\n✓ Product created"))
			fmt.Printf("  ID: %s\n", CyanStyle.Render(product.ID))
			fmt.Printf("  Identifier: %s\n", product.Identifier)
			fmt.Printf("  Type: %s\n", product.ProductType)
			return nil
		}

		fmt.Println(GreenStyle.Render("\n✓ Product created"))
		return nil
	})
}
func runProductsUpdate(cmd *cobra.Command, args []string) error {
	productID, _ := cmd.Flags().GetString("product-id")
	productType, _ := cmd.Flags().GetString("product-type")
	identifier, _ := cmd.Flags().GetString("identifier")
	displayName, _ := cmd.Flags().GetString("name")

	if productID == "" {
		return fmt.Errorf("--product-id is required")
	}

	patchBody := map[string]interface{}{}
	if productType != "" {
		if err := validateProductType(productType); err != nil {
			return err
		}
		patchBody["product_type"] = productType
	}
	if identifier != "" {
		patchBody["identifier"] = identifier
	}
	if displayName != "" {
		patchBody["display_name"] = displayName
	}
	if len(patchBody) == 0 {
		return fmt.Errorf("at least one of --product-type, --identifier, or --name is required")
	}

	c, err := Dashboard("\n🛍️ Updating product (PATCH)...")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/developers/me/projects/%s/products/%s", c.ProjectID, productID)
	resp, err := c.Client.Patch(path, patchBody)
	if err != nil {
		return err
	}
	return c.Respond(resp, func() error {

		fmt.Println(GreenStyle.Render("\n✓ Product updated via internal API."))
		return nil
	})
}
