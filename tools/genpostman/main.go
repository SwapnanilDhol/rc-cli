// Command genpostman writes postman/RevenueCat-API-v2.postman_collection.json
// Run from repo root: go run ./tools/genpostman
package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

// ep is one RevenueCat v2 route as of the published docs (Developer API v2).
type ep struct {
	Folder string
	Name   string
	Method string
	Path   string // starts with /; uses {{var}} for Postman
}

// internalEp is a dashboard session route.
// Base "": https://app.revenuecat.com/internal/v1 + Path (catalog, projects list, etc.)
// Base "appv1": https://app.revenuecat.com/v1 + Path (bare /developers/me and same-origin me routes — not under internal/v1).
// PostmanBody: optional default raw JSON in Postman for POST/PUT/PATCH; nil uses "{}".
type internalEp struct {
	Folder      string
	Name        string
	Method      string
	Path        string // e.g. /developers/me/projects — uses {{project_id}}, …
	Base        string // "" or "appv1"
	PostmanBody *string
}

func main() {
	repoRoot, err := findRepoRoot()
	if err != nil {
		panic(err)
	}
	outPath := filepath.Join(repoRoot, "postman", "RevenueCat-API-v2.postman_collection.json")

	endpoints := v2Endpoints()
	internalEndpoints := internalEndpointsFromCLI()

	folderOrder := []string{}
	seen := map[string]bool{}
	for _, e := range endpoints {
		if !seen[e.Folder] {
			seen[e.Folder] = true
			folderOrder = append(folderOrder, e.Folder)
		}
	}
	items := []interface{}{}
	for _, fname := range folderOrder {
		var sub []interface{}
		for _, e := range endpoints {
			if e.Folder != fname {
				continue
			}
			sub = append(sub, buildV2Item(e))
		}
		items = append(items, map[string]interface{}{
			"name": fname,
			"item": sub,
		})
	}

	items = append(items, buildInternalRootFolder(internalEndpoints))

	collection := map[string]interface{}{
		"info": map[string]interface{}{
			"name": "RevenueCat API — Developer v2 + Internal (dashboard)",
			"description": "Public v2: Bearer `apiKey` → `https://api.revenuecat.com/v2`. " +
				"Dashboard session: Cookie `rc_auth_token={{rc_auth_token}}`. Most internal JSON uses `internal/v1`; " +
				"**Account → GET developers/me** uses `https://app.revenuecat.com/v1` (not internal/v1) or you get HTTP 7117. " +
				"Run **Auth → Login** first. Official v2 docs: https://www.revenuecat.com/docs/api-v2. " +
				"Regenerate: `go run ./tools/genpostman` from repo root.",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"variable": collectionVariables(),
		"item":     items,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(collection); err != nil {
		panic(err)
	}
	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func collectionVariables() []interface{} {
	return []interface{}{
		map[string]interface{}{"key": "baseUrl", "value": "https://api.revenuecat.com/v2"},
		map[string]interface{}{"key": "apiKey", "value": ""},
		map[string]interface{}{"key": "internalBaseUrl", "value": "https://app.revenuecat.com/internal/v1"},
		map[string]interface{}{"key": "authBaseUrl", "value": "https://app.revenuecat.com/v1"},
		map[string]interface{}{"key": "rc_email", "value": ""},
		map[string]interface{}{"key": "rc_password", "value": ""},
		map[string]interface{}{"key": "rc_auth_token", "value": ""},
		map[string]interface{}{"key": "project_id", "value": ""},
		map[string]interface{}{"key": "customer_id", "value": ""},
		map[string]interface{}{"key": "app_id", "value": ""},
		map[string]interface{}{"key": "entitlement_id", "value": ""},
		map[string]interface{}{"key": "offering_id", "value": ""},
		map[string]interface{}{"key": "package_id", "value": ""},
		map[string]interface{}{"key": "product_id", "value": ""},
		map[string]interface{}{"key": "purchase_id", "value": ""},
		map[string]interface{}{"key": "subscription_id", "value": ""},
		map[string]interface{}{"key": "transaction_id", "value": ""},
		map[string]interface{}{"key": "paywall_id", "value": ""},
		map[string]interface{}{"key": "webhook_integration_id", "value": ""},
		map[string]interface{}{"key": "invoice_id", "value": ""},
		map[string]interface{}{"key": "virtual_currency_code", "value": ""},
		map[string]interface{}{"key": "chart_name", "value": "overview"},
		map[string]interface{}{"key": "experiment_id", "value": ""},
		map[string]interface{}{"key": "subscriber_list_id", "value": ""},
		map[string]interface{}{"key": "collaborator_id", "value": ""},
		map[string]interface{}{"key": "api_key_id", "value": ""},
		map[string]interface{}{"key": "webhook_id", "value": ""},
		map[string]interface{}{"key": "sandbox_mode", "value": "false"},
		map[string]interface{}{"key": "charts_resolution", "value": "86400"},
	}
}

func v2Endpoints() []ep {
	return []ep{
		// Projects
		{"Projects", "List projects", "GET", "/projects"},
		{"Projects", "Create project", "POST", "/projects"},
		{"Projects", "Get project", "GET", "/projects/{{project_id}}"},

		// Apps
		{"Apps", "List apps", "GET", "/projects/{{project_id}}/apps"},
		{"Apps", "Create app", "POST", "/projects/{{project_id}}/apps"},
		{"Apps", "Get app", "GET", "/projects/{{project_id}}/apps/{{app_id}}"},
		{"Apps", "Update app", "POST", "/projects/{{project_id}}/apps/{{app_id}}"},
		{"Apps", "Delete app", "DELETE", "/projects/{{project_id}}/apps/{{app_id}}"},
		{"Apps", "Get app public API keys", "GET", "/projects/{{project_id}}/apps/{{app_id}}/public_api_keys"},
		{"Apps", "Get StoreKit config", "GET", "/projects/{{project_id}}/apps/{{app_id}}/store_kit_config"},

		// Audit & collaborators
		{"Project admin", "List audit logs", "GET", "/projects/{{project_id}}/audit_logs"},
		{"Project admin", "List collaborators", "GET", "/projects/{{project_id}}/collaborators"},

		// Metrics & charts
		{"Metrics & charts", "Metrics overview", "GET", "/projects/{{project_id}}/metrics/overview"},
		{"Metrics & charts", "Get chart", "GET", "/projects/{{project_id}}/charts/{{chart_name}}"},
		{"Metrics & charts", "Get chart options", "GET", "/projects/{{project_id}}/charts/{{chart_name}}/options"},

		// Customers
		{"Customers", "List customers", "GET", "/projects/{{project_id}}/customers"},
		{"Customers", "Create customer", "POST", "/projects/{{project_id}}/customers"},
		{"Customers", "Get customer", "GET", "/projects/{{project_id}}/customers/{{customer_id}}"},
		{"Customers", "Delete customer", "DELETE", "/projects/{{project_id}}/customers/{{customer_id}}"},
		{"Customers", "Transfer customer", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/actions/transfer"},
		{"Customers", "Grant entitlement", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/actions/grant_entitlement"},
		{"Customers", "Revoke granted entitlement", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/actions/revoke_granted_entitlement"},
		{"Customers", "Assign offering override", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/actions/assign_offering"},
		{"Customers", "List subscriptions", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/subscriptions"},
		{"Customers", "List purchases", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/purchases"},
		{"Customers", "Active entitlements", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/active_entitlements"},
		{"Customers", "List aliases", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/aliases"},
		{"Customers", "Virtual currency balances", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/virtual_currencies"},
		{"Customers", "Virtual currency transaction", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/virtual_currencies/transactions"},
		{"Customers", "Virtual currency update balance", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/virtual_currencies/update_balance"},
		{"Customers", "Get attributes", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/attributes"},
		{"Customers", "Set attributes", "POST", "/projects/{{project_id}}/customers/{{customer_id}}/attributes"},
		{"Customers", "List invoices", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/invoices"},
		{"Customers", "Download invoice file", "GET", "/projects/{{project_id}}/customers/{{customer_id}}/invoices/{{invoice_id}}/file"},

		// Entitlements
		{"Entitlements", "List entitlements", "GET", "/projects/{{project_id}}/entitlements"},
		{"Entitlements", "Create entitlement", "POST", "/projects/{{project_id}}/entitlements"},
		{"Entitlements", "Get entitlement", "GET", "/projects/{{project_id}}/entitlements/{{entitlement_id}}"},
		{"Entitlements", "Update entitlement", "POST", "/projects/{{project_id}}/entitlements/{{entitlement_id}}"},
		{"Entitlements", "Delete entitlement", "DELETE", "/projects/{{project_id}}/entitlements/{{entitlement_id}}"},
		{"Entitlements", "List entitlement products", "GET", "/projects/{{project_id}}/entitlements/{{entitlement_id}}/products"},
		{"Entitlements", "Archive entitlement", "POST", "/projects/{{project_id}}/entitlements/{{entitlement_id}}/actions/archive"},
		{"Entitlements", "Unarchive entitlement", "POST", "/projects/{{project_id}}/entitlements/{{entitlement_id}}/actions/unarchive"},
		{"Entitlements", "Attach products", "POST", "/projects/{{project_id}}/entitlements/{{entitlement_id}}/actions/attach_products"},
		{"Entitlements", "Detach products", "POST", "/projects/{{project_id}}/entitlements/{{entitlement_id}}/actions/detach_products"},

		// Offerings
		{"Offerings", "List offerings", "GET", "/projects/{{project_id}}/offerings"},
		{"Offerings", "Create offering", "POST", "/projects/{{project_id}}/offerings"},
		{"Offerings", "Get offering", "GET", "/projects/{{project_id}}/offerings/{{offering_id}}"},
		{"Offerings", "Update offering", "POST", "/projects/{{project_id}}/offerings/{{offering_id}}"},
		{"Offerings", "Delete offering", "DELETE", "/projects/{{project_id}}/offerings/{{offering_id}}"},
		{"Offerings", "Archive offering", "POST", "/projects/{{project_id}}/offerings/{{offering_id}}/actions/archive"},
		{"Offerings", "Unarchive offering", "POST", "/projects/{{project_id}}/offerings/{{offering_id}}/actions/unarchive"},

		// Packages
		{"Packages", "List packages in offering", "GET", "/projects/{{project_id}}/offerings/{{offering_id}}/packages"},
		{"Packages", "Create package", "POST", "/projects/{{project_id}}/offerings/{{offering_id}}/packages"},
		{"Packages", "Get package", "GET", "/projects/{{project_id}}/packages/{{package_id}}"},
		{"Packages", "Update package", "POST", "/projects/{{project_id}}/packages/{{package_id}}"},
		{"Packages", "Delete package", "DELETE", "/projects/{{project_id}}/packages/{{package_id}}"},
		{"Packages", "List package products", "GET", "/projects/{{project_id}}/packages/{{package_id}}/products"},
		{"Packages", "Attach products", "POST", "/projects/{{project_id}}/packages/{{package_id}}/actions/attach_products"},
		{"Packages", "Detach products", "POST", "/projects/{{project_id}}/packages/{{package_id}}/actions/detach_products"},

		// Products
		{"Products", "List products", "GET", "/projects/{{project_id}}/products"},
		{"Products", "Create product", "POST", "/projects/{{project_id}}/products"},
		{"Products", "Get product", "GET", "/projects/{{project_id}}/products/{{product_id}}"},
		{"Products", "Update product", "POST", "/projects/{{project_id}}/products/{{product_id}}"},
		{"Products", "Delete product", "DELETE", "/projects/{{project_id}}/products/{{product_id}}"},
		{"Products", "Archive product", "POST", "/projects/{{project_id}}/products/{{product_id}}/actions/archive"},
		{"Products", "Unarchive product", "POST", "/projects/{{project_id}}/products/{{product_id}}/actions/unarchive"},
		{"Products", "Create in store", "POST", "/projects/{{project_id}}/products/{{product_id}}/create_in_store"},

		// Virtual currencies
		{"Virtual currencies", "List", "GET", "/projects/{{project_id}}/virtual_currencies"},
		{"Virtual currencies", "Create", "POST", "/projects/{{project_id}}/virtual_currencies"},
		{"Virtual currencies", "Get", "GET", "/projects/{{project_id}}/virtual_currencies/{{virtual_currency_code}}"},
		{"Virtual currencies", "Update", "POST", "/projects/{{project_id}}/virtual_currencies/{{virtual_currency_code}}"},
		{"Virtual currencies", "Delete", "DELETE", "/projects/{{project_id}}/virtual_currencies/{{virtual_currency_code}}"},
		{"Virtual currencies", "Archive", "POST", "/projects/{{project_id}}/virtual_currencies/{{virtual_currency_code}}/actions/archive"},
		{"Virtual currencies", "Unarchive", "POST", "/projects/{{project_id}}/virtual_currencies/{{virtual_currency_code}}/actions/unarchive"},

		// Purchases
		{"Purchases", "List purchases", "GET", "/projects/{{project_id}}/purchases"},
		{"Purchases", "Get purchase", "GET", "/projects/{{project_id}}/purchases/{{purchase_id}}"},
		{"Purchases", "Purchase entitlements", "GET", "/projects/{{project_id}}/purchases/{{purchase_id}}/entitlements"},
		{"Purchases", "Refund purchase", "POST", "/projects/{{project_id}}/purchases/{{purchase_id}}/actions/refund"},

		// Subscriptions
		{"Subscriptions", "List subscriptions", "GET", "/projects/{{project_id}}/subscriptions"},
		{"Subscriptions", "Get subscription", "GET", "/projects/{{project_id}}/subscriptions/{{subscription_id}}"},
		{"Subscriptions", "Subscription transactions", "GET", "/projects/{{project_id}}/subscriptions/{{subscription_id}}/transactions"},
		{"Subscriptions", "Refund subscription transaction", "POST", "/projects/{{project_id}}/subscriptions/{{subscription_id}}/transactions/{{transaction_id}}/actions/refund"},
		{"Subscriptions", "Subscription entitlements", "GET", "/projects/{{project_id}}/subscriptions/{{subscription_id}}/entitlements"},
		{"Subscriptions", "Cancel subscription", "POST", "/projects/{{project_id}}/subscriptions/{{subscription_id}}/actions/cancel"},
		{"Subscriptions", "Refund subscription", "POST", "/projects/{{project_id}}/subscriptions/{{subscription_id}}/actions/refund"},
		{"Subscriptions", "Authenticated management URL", "GET", "/projects/{{project_id}}/subscriptions/{{subscription_id}}/authenticated_management_url"},

		// Paywalls
		{"Paywalls", "List paywalls", "GET", "/projects/{{project_id}}/paywalls"},
		{"Paywalls", "Create paywall", "POST", "/projects/{{project_id}}/paywalls"},
		{"Paywalls", "Get paywall", "GET", "/projects/{{project_id}}/paywalls/{{paywall_id}}"},
		{"Paywalls", "Delete paywall", "DELETE", "/projects/{{project_id}}/paywalls/{{paywall_id}}"},

		// Webhook integrations (v2)
		{"Integrations", "List webhook integrations", "GET", "/projects/{{project_id}}/integrations/webhooks"},
		{"Integrations", "Create webhook integration", "POST", "/projects/{{project_id}}/integrations/webhooks"},
		{"Integrations", "Get webhook integration", "GET", "/projects/{{project_id}}/integrations/webhooks/{{webhook_integration_id}}"},
		{"Integrations", "Update webhook integration", "POST", "/projects/{{project_id}}/integrations/webhooks/{{webhook_integration_id}}"},
		{"Integrations", "Delete webhook integration", "DELETE", "/projects/{{project_id}}/integrations/webhooks/{{webhook_integration_id}}"},
	}
}

func internalEndpointsFromCLI() []internalEp {
	postmanOfferingMetadata := "{\n  \"display_name\": \"Edit in Postman\",\n  \"metadata\": {\n    \"rc_cli_note\": \"dashboard should show this under offering / Paywalls context\"\n  }\n}"
	postmanProductMutationBody := "{\n  \"product_type\": \"subscription\",\n  \"identifier\": \"rc.identifier\",\n  \"display_name\": \"Display name\"\n}"
	postmanAppStoreProductsCreateBody := "{\n  \"products\": [\n    {\n      \"product_identifier\": \"rc.identifier\",\n      \"name\": \"Display name\",\n      \"product_type\": \"subscriptions\",\n      \"duration\": \"ONE_WEEK\",\n      \"subscription_group\": {\n        \"id\": \"21492027\",\n        \"name\": \"passMaker-subscription-group\"\n      }\n    }\n  ]\n}"
	return []internalEp{
		// Account / me — same-origin v1 (NOT internal/v1); see API.md § Same-origin v1
		{"Account", "GET developers/me", "GET", "/developers/me", "appv1", nil},
		{"Account", "GET billing info", "GET", "/developers/me/billing/info", "appv1", nil},
		{"Account", "GET pending collaborations", "GET", "/developers/me/collaborations/pending", "appv1", nil},
		{"Account", "GET dashboard notifications", "GET", "/developers/me/dashboard_notifications", "appv1", nil},
		{"Account", "GET transactions (me)", "GET", "/developers/me/transactions", "appv1", nil},

		// Projects
		{"Internal — Projects", "List all projects", "GET", "/developers/me/projects", "", nil},
		{"Internal — Projects", "Create project", "POST", "/developers/me/projects", "", nil},
		{"Internal — Projects", "Get project", "GET", "/developers/me/projects/{{project_id}}", "", nil},

		// Catalog
		{"Internal — Catalog", "List entitlements", "GET", "/developers/me/projects/{{project_id}}/entitlements", "", nil},
		{"Internal — Catalog", "Create entitlement", "POST", "/developers/me/projects/{{project_id}}/entitlements", "", nil},
		{"Internal — Catalog", "Get entitlement", "GET", "/developers/me/projects/{{project_id}}/entitlements/{{entitlement_id}}", "", nil},
		{"Internal — Catalog", "Update entitlement", "PUT", "/developers/me/projects/{{project_id}}/entitlements/{{entitlement_id}}", "", nil},
		{"Internal — Catalog", "Delete entitlement", "DELETE", "/developers/me/projects/{{project_id}}/entitlements/{{entitlement_id}}", "", nil},
		{"Internal — Catalog", "Archive entitlement", "POST", "/developers/me/projects/{{project_id}}/entitlements/{{entitlement_id}}/actions/archive", "", nil},
		{"Internal — Catalog", "Get entitlement products", "GET", "/developers/me/projects/{{project_id}}/entitlements/{{entitlement_id}}/products", "", nil},
		{"Internal — Catalog", "List offerings", "GET", "/developers/me/projects/{{project_id}}/offerings", "", nil},
		{"Internal — Catalog", "Create offering", "POST", "/developers/me/projects/{{project_id}}/offerings", "", nil},
		{"Internal — Catalog", "Get offering", "GET", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}", "", nil},
		{"Internal — Catalog", "Update offering", "PUT", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}", "", nil},
		{"Internal — Catalog", "Update offering (sample body: metadata)", "PUT", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}", "", &postmanOfferingMetadata},
		{"Internal — Catalog", "Set current offering (PATCH)", "PATCH", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}", "", nil},
		{"Internal — Catalog", "Duplicate offering", "POST", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}/duplicate", "", nil},
		{"Internal — Catalog", "Archive offering", "POST", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}/actions/archive", "", nil},
		{"Internal — Catalog", "Delete offering", "DELETE", "/developers/me/projects/{{project_id}}/offerings/{{offering_id}}", "", nil},
		{"Internal — Catalog", "List products", "GET", "/developers/me/projects/{{project_id}}/products", "", nil},
		{"Internal — Products", "Create product for app", "POST", "/developers/me/projects/{{project_id}}/apps/{{app_id}}/products", "", &postmanProductMutationBody},
		{"Internal — Products", "Update product (PATCH)", "PATCH", "/developers/me/projects/{{project_id}}/products/{{product_id}}", "", &postmanProductMutationBody},
		{"Internal — Catalog", "List paywalls", "GET", "/developers/me/projects/{{project_id}}/paywalls", "", nil},
		{"Internal — Catalog", "Product store statuses", "GET", "/developers/me/projects/{{project_id}}/product_stores_statuses", "", nil},
		{"Internal — Catalog", "List promotions", "GET", "/developers/me/projects/{{project_id}}/promotions", "", nil},
		{"Internal — Catalog", "List intro offers", "GET", "/developers/me/projects/{{project_id}}/intro_offers", "", nil},
		{"Internal — Catalog", "List apps", "GET", "/developers/me/projects/{{project_id}}/apps", "", nil},
		{"Internal — Catalog", "List subscription groups (App Store Connect)", "GET", "/developers/me/projects/{{project_id}}/apps/{{app_id}}/subscription_groups", "", nil},
		{"Internal — Catalog", "Create app store product", "POST", "/developers/me/projects/{{project_id}}/apps/{{app_id}}/app_store_products", "", &postmanAppStoreProductsCreateBody},

		// Project admin
		{"Internal — Project admin", "List collaborators", "GET", "/developers/me/projects/{{project_id}}/collaborators", "", nil},
		{"Internal — Project admin", "Get collaborator", "GET", "/developers/me/projects/{{project_id}}/collaborators/{{collaborator_id}}", "", nil},
		{"Internal — Project admin", "Add collaborator", "POST", "/developers/me/projects/{{project_id}}/collaborators", "", nil},
		{"Internal — Project admin", "Update collaborator", "PUT", "/developers/me/projects/{{project_id}}/collaborators/{{collaborator_id}}", "", nil},
		{"Internal — Project admin", "Remove collaborator", "DELETE", "/developers/me/projects/{{project_id}}/collaborators/{{collaborator_id}}", "", nil},
		{"Internal — Project admin", "List API keys", "GET", "/developers/me/projects/{{project_id}}/api_keys", "", nil},
		{"Internal — Project admin", "Create API key", "POST", "/developers/me/projects/{{project_id}}/api_keys", "", nil},
		{"Internal — Project admin", "Delete API key", "DELETE", "/developers/me/projects/{{project_id}}/api_keys/{{api_key_id}}", "", nil},
		{"Internal — Project admin", "List audit logs", "GET", "/developers/me/projects/{{project_id}}/audit_logs", "", nil},
		{"Internal — Project admin", "List webhooks (legacy path)", "GET", "/developers/me/projects/{{project_id}}/webhooks", "", nil},
		{"Internal — Project admin", "Create webhook", "POST", "/developers/me/projects/{{project_id}}/webhooks", "", nil},
		{"Internal — Project admin", "Update webhook", "PUT", "/developers/me/projects/{{project_id}}/webhooks/{{webhook_id}}", "", nil},
		{"Internal — Project admin", "Delete webhook", "DELETE", "/developers/me/projects/{{project_id}}/webhooks/{{webhook_id}}", "", nil},
		{"Internal — Project admin", "Test webhook", "POST", "/developers/me/projects/{{project_id}}/webhooks/{{webhook_id}}/test", "", nil},
		{"Internal — Project admin", "Webhook events", "GET", "/developers/me/projects/{{project_id}}/webhooks/events", "", nil},

		// Experiments & lists
		{"Internal — Experiments & lists", "List price experiments", "GET", "/developers/me/projects/{{project_id}}/price_experiments", "", nil},
		{"Internal — Experiments & lists", "Get price experiment", "GET", "/developers/me/projects/{{project_id}}/price_experiments/{{experiment_id}}", "", nil},
		{"Internal — Experiments & lists", "Create price experiment", "POST", "/developers/me/projects/{{project_id}}/price_experiments", "", nil},
		{"Internal — Experiments & lists", "Pause price experiment", "POST", "/developers/me/projects/{{project_id}}/price_experiments/{{experiment_id}}/pause", "", nil},
		{"Internal — Experiments & lists", "Resume price experiment", "POST", "/developers/me/projects/{{project_id}}/price_experiments/{{experiment_id}}/resume", "", nil},
		{"Internal — Experiments & lists", "Stop price experiment", "POST", "/developers/me/projects/{{project_id}}/price_experiments/{{experiment_id}}/stop", "", nil},
		{"Internal — Experiments & lists", "List subscriber lists", "GET", "/developers/me/projects/{{project_id}}/subscriber_lists", "", nil},
		{"Internal — Experiments & lists", "Get subscriber list", "GET", "/developers/me/subscriber_lists/{{subscriber_list_id}}", "", nil},
		{"Internal — Experiments & lists", "Subscriber lists manifest", "GET", "/developers/me/subscriber_lists/manifest", "", nil},

		// Charts (dashboard analytics; query params match rc internal charts)
		{"Internal — Charts v2", "Overview (project)", "GET", "/developers/me/charts_v2/overview?app_uuid={{project_id}}&sandbox_mode={{sandbox_mode}}", "", nil},
		{"Internal — Charts v2", "Overview (all projects)", "GET", "/developers/me/charts_v2/overview?sandbox_mode={{sandbox_mode}}&v3=false", "", nil},
		{"Internal — Charts v2", "Trials", "GET", "/developers/me/charts_v2/trials?app_uuid={{project_id}}&sandbox_mode={{sandbox_mode}}&resolution={{charts_resolution}}", "", nil},
		{"Internal — Charts v2", "Transactions", "GET", "/developers/me/charts_v2/transactions?app_uuid={{project_id}}&sandbox_mode={{sandbox_mode}}&resolution={{charts_resolution}}", "", nil},
		{"Internal — Charts v2", "Revenue", "GET", "/developers/me/charts_v2/revenue?app_uuid={{project_id}}&sandbox_mode={{sandbox_mode}}&resolution={{charts_resolution}}", "", nil},
		{"Internal — Charts v2", "Customers new", "GET", "/developers/me/charts_v2/customers_new?sandbox_mode={{sandbox_mode}}", "", nil},
		{"Internal — Charts v2", "MRR", "GET", "/developers/me/charts_v2/mrr?sandbox_mode={{sandbox_mode}}", "", nil},
		{"Internal — Charts v2", "Actives", "GET", "/developers/me/charts_v2/actives?sandbox_mode={{sandbox_mode}}", "", nil},
	}
}

func buildV2Item(e ep) map[string]interface{} {
	rawURL := "{{baseUrl}}" + e.Path
	req := map[string]interface{}{
		"method": e.Method,
		"header": []map[string]string{
			{"key": "Authorization", "value": "Bearer {{apiKey}}"},
		},
		"url": rawURL,
	}
	if e.Method == "POST" || e.Method == "PUT" || e.Method == "PATCH" {
		req["header"] = []map[string]string{
			{"key": "Authorization", "value": "Bearer {{apiKey}}"},
			{"key": "Content-Type", "value": "application/json"},
		}
		req["body"] = map[string]interface{}{
			"mode": "raw",
			"raw":  "{}",
		}
	}
	return map[string]interface{}{
		"name":    e.Name,
		"request": req,
	}
}

func internalHeaders(method string) []map[string]string {
	h := []map[string]string{
		{"key": "Accept", "value": "application/json"},
		{"key": "Content-Type", "value": "application/json"},
		{"key": "X-Requested-With", "value": "XMLHttpRequest"},
		{"key": "Cookie", "value": "rc_auth_token={{rc_auth_token}}"},
	}
	if method != "GET" && method != "DELETE" {
		h = append(h,
			map[string]string{"key": "Origin", "value": "https://app.revenuecat.com"},
			map[string]string{"key": "Referer", "value": "https://app.revenuecat.com/"},
		)
	}
	return h
}

func buildInternalItem(e internalEp) map[string]interface{} {
	prefix := "{{internalBaseUrl}}"
	if e.Base == "appv1" {
		prefix = "{{authBaseUrl}}"
	}
	rawURL := prefix + e.Path
	req := map[string]interface{}{
		"method": e.Method,
		"header": internalHeaders(e.Method),
		"url":    rawURL,
	}
	if e.Method == "POST" || e.Method == "PUT" || e.Method == "PATCH" {
		raw := "{}"
		if e.PostmanBody != nil {
			raw = *e.PostmanBody
		}
		req["body"] = map[string]interface{}{
			"mode": "raw",
			"raw":  raw,
		}
	}
	return map[string]interface{}{
		"name":    e.Name,
		"request": req,
	}
}

func buildAuthLoginItem() map[string]interface{} {
	return map[string]interface{}{
		"name": "Login (email + password → sets rc_auth_token)",
		"request": map[string]interface{}{
			"method": "POST",
			"header": []map[string]string{
				{"key": "Accept", "value": "application/json"},
				{"key": "Content-Type", "value": "application/json"},
				{"key": "X-Requested-With", "value": "XMLHttpRequest"},
				{"key": "Origin", "value": "https://app.revenuecat.com"},
				{"key": "Referer", "value": "https://app.revenuecat.com/login"},
			},
			"body": map[string]interface{}{
				"mode": "raw",
				"raw":  "{\n  \"email\": \"{{rc_email}}\",\n  \"password\": \"{{rc_password}}\"\n}",
			},
			"url":         "{{authBaseUrl}}/developers/login",
			"description": "Same as `rc login`. Copies `authentication_token` into collection variable `rc_auth_token` (Tests tab).",
		},
		"event": []interface{}{
			map[string]interface{}{
				"listen": "test",
				"script": map[string]interface{}{
					"type": "text/javascript",
					"exec": []string{
						"try {",
						"    var j = pm.response.json();",
						"    if (j && j.authentication_token) {",
						"        pm.collectionVariables.set(\"rc_auth_token\", j.authentication_token);",
						"        console.log(\"Set rc_auth_token\");",
						"    }",
						"} catch (e) { console.warn(e); }",
					},
				},
			},
		},
	}
}

func buildAuthRefreshItem() map[string]interface{} {
	return map[string]interface{}{
		"name": "Refresh session token",
		"request": map[string]interface{}{
			"method": "POST",
			"header": []map[string]string{
				{"key": "Accept", "value": "application/json"},
				{"key": "Content-Type", "value": "application/json"},
				{"key": "X-Requested-With", "value": "XMLHttpRequest"},
				{"key": "Origin", "value": "https://app.revenuecat.com"},
				{"key": "Referer", "value": "https://app.revenuecat.com/login"},
				{"key": "Cookie", "value": "rc_auth_token={{rc_auth_token}}"},
			},
			"body": map[string]interface{}{
				"mode": "raw",
				"raw":  "{}",
			},
			"url":         "{{authBaseUrl}}/developers/login/refresh-token",
			"description": "Same as CLI refresh. Requires existing `rc_auth_token`. Updates variable from response.",
		},
		"event": []interface{}{
			map[string]interface{}{
				"listen": "test",
				"script": map[string]interface{}{
					"type": "text/javascript",
					"exec": []string{
						"try {",
						"    var j = pm.response.json();",
						"    if (j && j.authentication_token) {",
						"        pm.collectionVariables.set(\"rc_auth_token\", j.authentication_token);",
						"    }",
						"} catch (e) { console.warn(e); }",
					},
				},
			},
		},
	}
}

func buildInternalRootFolder(endpoints []internalEp) map[string]interface{} {
	authFolder := map[string]interface{}{
		"name":        "Auth",
		"description": "Set `rc_email` and `rc_password` in collection variables, run **Login** — then use Internal requests.",
		"item": []interface{}{
			buildAuthLoginItem(),
			buildAuthRefreshItem(),
		},
	}

	folderOrder := []string{}
	seen := map[string]bool{}
	for _, e := range endpoints {
		if !seen[e.Folder] {
			seen[e.Folder] = true
			folderOrder = append(folderOrder, e.Folder)
		}
	}

	rest := []interface{}{}
	for _, fname := range folderOrder {
		var sub []interface{}
		for _, e := range endpoints {
			if e.Folder != fname {
				continue
			}
			sub = append(sub, buildInternalItem(e))
		}
		rest = append(rest, map[string]interface{}{
			"name": fname,
			"item": sub,
		})
	}

	all := []interface{}{authFolder}
	all = append(all, rest...)

	return map[string]interface{}{
		"name": "Internal — dashboard session",
		"description": "Undocumented dashboard JSON. Auth: Cookie `rc_auth_token` after **Auth → Login**. " +
			"Most routes: `{{internalBaseUrl}}` …/internal/v1. **Account** folder: `{{authBaseUrl}}` …/v1 (e.g. GET /developers/me is not served under internal/v1). Matches `revenuecat-cli` (see `internal/internal.go`).",
		"item": all,
	}
}
