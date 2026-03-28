package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/config"
	"revenuecat-cli/internal"
)

var internalStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
var greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
var yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var grayStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

func init() {
	// Projects command
	projectsCmd := &cobra.Command{
		Use:   "projects",
		Short: "Manage projects",
	}
	RootCmd.AddCommand(projectsCmd)

	projectsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		RunE:  runInternalProjectsList,
	}
	projectsListCmd.Aliases = []string{"ls"}

	projectsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get project details",
		RunE:  runInternalProjectGet,
	}
	projectsGetCmd.Flags().StringP("project-id", "i", "", "Project ID")

	projectsCmd.AddCommand(projectsListCmd, projectsGetCmd)

	projectsUseCmd := &cobra.Command{
		Use:   "use",
		Short: "Set default project",
		RunE:  runInternalProjectsUse,
	}
	projectsUseCmd.Flags().StringP("project-id", "i", "", "Project ID (required)")

	projectsCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project",
		RunE:  runInternalProjectsCreate,
	}
	projectsCreateCmd.Flags().StringP("name", "n", "", "Project name (required)")
	projectsCreateCmd.Flags().String("category", "general", "Project category")
	projectsCreateCmd.Flags().StringSlice("platforms", []string{"apple_native"}, "Expected platforms (apple_native, google_play, stripe, amazon, etc.)")

	projectsCmd.AddCommand(projectsUseCmd, projectsCreateCmd)

	// Entitlements command
	entitlementsCmd := &cobra.Command{
		Use:   "entitlements",
		Short: "Manage entitlements",
	}
	RootCmd.AddCommand(entitlementsCmd)

	entitlementsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all entitlements",
		RunE:  runInternalEntitlementsList,
	}
	entitlementsListCmd.Aliases = []string{"ls"}

	entitlementsCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an entitlement",
		RunE:  runInternalEntitlementsCreate,
	}
	entitlementsCreateCmd.Flags().StringP("identifier", "i", "", "Entitlement identifier (slug)")
	entitlementsCreateCmd.Flags().StringP("name", "n", "", "Display name")

	entitlementsDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an entitlement",
		RunE:  runInternalEntitlementsDelete,
	}
	entitlementsDeleteCmd.Flags().StringP("entitlement-id", "e", "", "Entitlement ID")

	entitlementsCmd.AddCommand(entitlementsListCmd, entitlementsCreateCmd, entitlementsDeleteCmd)

	// Offerings command
	offeringsCmd := &cobra.Command{
		Use:   "offerings",
		Short: "Manage offerings",
	}
	RootCmd.AddCommand(offeringsCmd)

	offeringsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all offerings",
		RunE:  runInternalOfferingsList,
	}
	offeringsListCmd.Flags().String("platform", "", "Filter by platform (IOS, ANDROID, MACOS, WINDOWS, LINUX, tvOS)")
	offeringsListCmd.Aliases = []string{"ls"}

	offeringsCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an offering",
		RunE:  runInternalOfferingsCreate,
	}
	offeringsCreateCmd.Flags().StringP("identifier", "i", "", "Offering identifier (slug)")
	offeringsCreateCmd.Flags().StringP("name", "n", "", "Display name")

	offeringsDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an offering",
		RunE:  runInternalOfferingsDelete,
	}
	offeringsDeleteCmd.Flags().StringP("offering-id", "o", "", "Offering ID")

	offeringsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get offering details",
		RunE:  runInternalOfferingGet,
	}
	offeringsGetCmd.Flags().StringP("offering-id", "o", "", "Offering ID")

	offeringsDuplicateCmd := &cobra.Command{
		Use:   "duplicate",
		Short: "Duplicate an offering",
		RunE:  runInternalOfferingsDuplicate,
	}
	offeringsDuplicateCmd.Flags().StringP("offering-id", "o", "", "Offering ID to duplicate (required)")
	offeringsDuplicateCmd.Flags().StringP("identifier", "i", "", "New offering identifier (required)")
	offeringsDuplicateCmd.Flags().StringP("name", "n", "", "New offering display name (required)")
	offeringsDuplicateCmd.Flags().Bool("packages-only", false, "Only duplicate packages")

	offeringsSetCurrentCmd := &cobra.Command{
		Use:   "set-current",
		Short: "Set an offering as the current/default",
		RunE:  runInternalOfferingsSetCurrent,
	}
	offeringsSetCurrentCmd.Flags().StringP("offering-id", "o", "", "Offering ID (required)")

	offeringsCmd.AddCommand(offeringsListCmd, offeringsGetCmd, offeringsCreateCmd, offeringsDeleteCmd, offeringsDuplicateCmd, offeringsSetCurrentCmd)

	// Products command
	productsCmd := &cobra.Command{
		Use:   "products",
		Short: "Manage products",
	}
	RootCmd.AddCommand(productsCmd)

	productsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all products",
		RunE:  runInternalProductsList,
	}
	productsListCmd.Flags().Int("limit", 100, "Maximum number of products to fetch")
	productsListCmd.Aliases = []string{"ls"}

	productsCmd.AddCommand(productsListCmd)

	// Apps command
	appsCmd := &cobra.Command{
		Use:   "apps",
		Short: "Manage apps",
	}
	RootCmd.AddCommand(appsCmd)

	appsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all apps",
		RunE:  runInternalAppsList,
	}
	appsListCmd.Aliases = []string{"ls"}

	appsCmd.AddCommand(appsListCmd)

	// Product Stores Status command
	storesStatusCmd := &cobra.Command{
		Use:   "stores-status",
		Short: "Get product stores connection status",
		RunE:  runInternalProductStoresStatuses,
	}
	storesStatusCmd.Aliases = []string{"stores", "product-stores"}
	RootCmd.AddCommand(storesStatusCmd)

	// Collaborators command
	collaboratorsCmd := &cobra.Command{
		Use:   "collaborators",
		Short: "Manage collaborators",
	}
	RootCmd.AddCommand(collaboratorsCmd)

	collaboratorsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List collaborators",
		RunE:  runInternalCollaboratorsList,
	}
	collaboratorsListCmd.Aliases = []string{"ls"}

	collaboratorsCmd.AddCommand(collaboratorsListCmd)

	// API Keys command
	apiKeysCmd := &cobra.Command{
		Use:   "apikeys",
		Short: "Manage API keys",
	}
	RootCmd.AddCommand(apiKeysCmd)

	apiKeysListCmd := &cobra.Command{
		Use:   "list",
		Short: "List API keys",
		RunE:  runInternalAPIKeysList,
	}
	apiKeysListCmd.Aliases = []string{"ls"}

	apiKeysCmd.AddCommand(apiKeysListCmd)

	// Audit logs command
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "View audit logs",
	}
	RootCmd.AddCommand(auditCmd)

	auditListCmd := &cobra.Command{
		Use:   "list",
		Short: "List audit logs",
		RunE:  runInternalAuditList,
	}
	auditListCmd.Aliases = []string{"ls"}

	auditCmd.AddCommand(auditListCmd)

	// Price Experiments command
	experimentsCmd := &cobra.Command{
		Use:   "experiments",
		Short: "Manage price experiments (A/B tests)",
	}
	RootCmd.AddCommand(experimentsCmd)

	experimentsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List price experiments",
		RunE:  runInternalPriceExperimentsList,
	}
	experimentsListCmd.Aliases = []string{"ls"}

	experimentsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get price experiment details",
		RunE:  runInternalPriceExperimentGet,
	}
	experimentsGetCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID")

	experimentsCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a price experiment",
		RunE:  runInternalPriceExperimentCreate,
	}
	experimentsCreateCmd.Flags().StringP("name", "n", "", "Experiment display name (required)")
	experimentsCreateCmd.Flags().StringP("offering-a", "a", "", "Offering A ID (required)")
	experimentsCreateCmd.Flags().StringP("offering-b", "b", "", "Offering B ID (required)")
	experimentsCreateCmd.Flags().Int("enrollment", 100, "Enrollment percentage (default 100)")
	experimentsCreateCmd.Flags().String("type", "introductory_offer", "Experiment type")
	experimentsCreateCmd.Flags().String("primary-metric", "Realized LTV per customer", "Primary metric")
	experimentsCreateCmd.Flags().String("notes", "", "Experiment notes")

	experimentsPauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause a running experiment",
		RunE:  runInternalPriceExperimentPause,
	}
	experimentsPauseCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID (required)")

	experimentsResumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume a paused experiment",
		RunE:  runInternalPriceExperimentResume,
	}
	experimentsResumeCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID (required)")

	experimentsStopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an experiment",
		RunE:  runInternalPriceExperimentStop,
	}
	experimentsStopCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID (required)")

	experimentsCmd.AddCommand(experimentsListCmd, experimentsGetCmd, experimentsCreateCmd, experimentsPauseCmd, experimentsResumeCmd, experimentsStopCmd)

	// Utilities command
	utilitiesCmd := &cobra.Command{
		Use:   "utilities",
		Short: "Utility endpoints",
	}
	RootCmd.AddCommand(utilitiesCmd)

	utilitiesCountriesCmd := &cobra.Command{
		Use:   "countries",
		Short: "List supported countries",
		RunE:  runInternalUtilitiesCountries,
	}

	utilitiesCmd.AddCommand(utilitiesCountriesCmd)

	// Subscriber Lists command
	listsCmd := &cobra.Command{
		Use:   "lists",
		Short: "Manage subscriber lists",
	}
	RootCmd.AddCommand(listsCmd)

	listsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List subscriber lists",
		RunE:  runInternalSubscriberListsList,
	}
	listsListCmd.Flags().Int("limit", 100, "Maximum number of lists to fetch")
	listsListCmd.Aliases = []string{"ls"}

	listsGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get subscriber list details",
		RunE:  runInternalSubscriberListGet,
	}
	listsGetCmd.Flags().StringP("list-id", "l", "", "List ID (required)")

	listsManifestCmd := &cobra.Command{
		Use:   "manifest",
		Short: "Get manifest of all subscriber lists",
		RunE:  runInternalSubscriberListManifest,
	}

	listsCmd.AddCommand(listsListCmd, listsGetCmd, listsManifestCmd)

	// Charts V2 command
	chartsCmd := &cobra.Command{
		Use:   "charts",
		Short: "View analytics charts",
	}
	RootCmd.AddCommand(chartsCmd)

	chartsOverviewCmd := &cobra.Command{
		Use:   "overview",
		Short: "Get project overview analytics",
		RunE:  runInternalChartsOverview,
	}
	chartsOverviewCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsOverviewCmd.Flags().String("app-uuid", "", "App UUID (optional, uses first app if not specified)")

	chartsTrialsCmd := &cobra.Command{
		Use:   "trials",
		Short: "Get trial analytics",
		RunE:  runInternalChartsTrials,
	}
	chartsTrialsCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	chartsTrialsCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	chartsTrialsCmd.Flags().Int("resolution", 0, "Resolution (0=daily, 1=weekly, 2=monthly)")
	chartsTrialsCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsTrialsCmd.Flags().String("app-uuid", "", "App UUID")

	chartsOverviewAllCmd := &cobra.Command{
		Use:   "overview-all",
		Short: "Get overview analytics across all projects",
		RunE:  runInternalChartsOverviewAll,
	}
	chartsOverviewAllCmd.Flags().Bool("sandbox", false, "Include sandbox data")

	chartsTransactionsCmd := &cobra.Command{
		Use:   "transactions",
		Short: "Get transaction analytics",
		RunE:  runInternalChartsTransactions,
	}
	chartsTransactionsCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	chartsTransactionsCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	chartsTransactionsCmd.Flags().Int("resolution", 0, "Resolution (0=daily, 1=weekly, 2=monthly)")
	chartsTransactionsCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsTransactionsCmd.Flags().String("app-uuid", "", "App UUID")

	chartsRevenueCmd := &cobra.Command{
		Use:   "revenue",
		Short: "Get revenue analytics",
		RunE:  runInternalChartsRevenue,
	}
	chartsRevenueCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	chartsRevenueCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	chartsRevenueCmd.Flags().Int("resolution", 0, "Resolution (0=daily, 1=weekly, 2=monthly)")
	chartsRevenueCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsRevenueCmd.Flags().String("app-uuid", "", "App UUID")

	chartsCmd.AddCommand(chartsOverviewCmd, chartsOverviewAllCmd, chartsTrialsCmd, chartsTransactionsCmd, chartsRevenueCmd)
}

func getInternalClient() (*internal.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	client, err := internal.EnsureAuthenticated(cfg)
	if err != nil {
		return nil, err
	}
	// If token was refreshed, save the new token
	if cfg.AuthToken != "" {
		if saveErr := config.SaveConfig(cfg); saveErr != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", saveErr)
		}
	}
	return client, nil
}

func getProjectID() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	if cfg.ProjectID == "" {
		return "", fmt.Errorf("no project selected. Run: rc projects list && rc projects use <id>")
	}
	return cfg.ProjectID, nil
}

func runInternalProjectsList(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📁 Fetching projects...")

	path := "/developers/me/projects"
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var projects []internal.Project
	if err := json.Unmarshal(toJSON(resp.Items), &projects); err != nil {
		return fmt.Errorf("error parsing projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println(yellowStyle.Render("No projects found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n📁 Your Projects:\n"))
	for _, p := range projects {
		currentMarker := ""
		if p.ID == cfg.ProjectID {
			currentMarker = greenStyle.Render(" (current)")
		}
		fmt.Printf("  ID: %s%s\n", cyanStyle.Render(p.ID), currentMarker)
		fmt.Printf("  Name: %s\n", p.Name)
		fmt.Printf("  Owner: %s\n", p.OwnerEmail)
		if p.RestrictedAccess {
			fmt.Println("  ⚠ Restricted access")
		}
		fmt.Println()
	}
	fmt.Println(grayStyle.Render(fmt.Sprintf("Total: %d projects", len(projects))))
	fmt.Println(grayStyle.Render("\nUse 'rc projects use <id>' to set default project"))

	return nil
}

func runInternalProjectGet(cmd *cobra.Command, args []string) error {
	projectID, _ := cmd.Flags().GetString("project-id")
	if projectID == "" {
		return fmt.Errorf("project-id is required")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📁 Fetching project details...")

	path := fmt.Sprintf("/developers/me/projects/%s", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	// Parse as generic map
	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📁 Project Details:\n"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(projectID))
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
}

func runInternalProjectsUse(cmd *cobra.Command, args []string) error {
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

	fmt.Println(greenStyle.Render("\n✓ Default project set to: ") + cyanStyle.Render(projectID))
	return nil
}

func runInternalProjectsCreate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}

	category, _ := cmd.Flags().GetString("category")
	platforms, _ := cmd.Flags().GetStringSlice("platforms")

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📁 Creating project...")

	path := "/developers/me/projects"
	data := map[string]interface{}{
		"name":                name,
		"category":           category,
		"expected_platforms": platforms,
	}

	resp, err := client.Post(path, data)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	var project map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &project); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(greenStyle.Render("\n✓ Project created:"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(project["id"].(string)))
	fmt.Printf("  Name: %s\n", project["name"])

	// Auto-set as default project
	cfg, err := config.LoadConfig()
	if err == nil {
		cfg.ProjectID = project["id"].(string)
		config.SaveConfig(cfg)
		fmt.Println(greenStyle.Render("  Set as default project"))
	}

	return nil
}

func runInternalEntitlementsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Fetching entitlements...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var entitlements []internal.Entitlement
	if err := json.Unmarshal(toJSON(resp.Items), &entitlements); err != nil {
		return fmt.Errorf("error parsing entitlements: %w", err)
	}

	if len(entitlements) == 0 {
		fmt.Println(yellowStyle.Render("No entitlements found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n📋 Entitlements:\n"))
	for _, e := range entitlements {
		fmt.Printf("  ID: %s\n", cyanStyle.Render(e.ID))
		fmt.Printf("  Identifier: %s\n", e.Identifier)
		fmt.Printf("  Name: %s\n", e.DisplayName)
		fmt.Printf("  Products: %d\n", len(e.Products))
		fmt.Println()
	}
	fmt.Println(grayStyle.Render(fmt.Sprintf("Total: %d entitlements", len(entitlements))))

	return nil
}

func runInternalEntitlementsCreate(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
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

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Creating entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements", projectID)
	data := map[string]string{
		"identifier":   identifier,
		"display_name": name,
	}
	resp, err := client.Post(path, data)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	var entitlement internal.Entitlement
	if err := json.Unmarshal(toJSON(resp.Data), &entitlement); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(greenStyle.Render("\n✓ Entitlement created:"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(entitlement.ID))
	fmt.Printf("  Identifier: %s\n", entitlement.Identifier)
	fmt.Printf("  Name: %s\n", entitlement.DisplayName)

	return nil
}

func runInternalEntitlementsDelete(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	entitlementID, _ := cmd.Flags().GetString("entitlement-id")
	if entitlementID == "" {
		return fmt.Errorf("entitlement-id is required (--entitlement-id or -e)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Deleting entitlement...")

	path := fmt.Sprintf("/developers/me/projects/%s/entitlements/%s", projectID, entitlementID)
	_, err = client.Delete(path)
	if err != nil {
		return err
	}

	fmt.Println(greenStyle.Render("\n✓ Entitlement deleted: " + entitlementID))
	return nil
}

func runInternalOfferingsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	platform, _ := cmd.Flags().GetString("platform")

	fmt.Println("\n📦 Fetching offerings...")

	path := fmt.Sprintf("/developers/me/projects/%s/offerings", projectID)
	params := map[string]string{}
	if platform != "" {
		params["platform"] = platform
	}
	var resp *internal.Response
	if len(params) > 0 {
		resp, err = client.GetWithParams(path, params)
	} else {
		resp, err = client.Get(path)
	}
	if err != nil {
		return err
	}

	var offerings []internal.Offering
	if err := json.Unmarshal(toJSON(resp.Items), &offerings); err != nil {
		return fmt.Errorf("error parsing offerings: %w", err)
	}

	if len(offerings) == 0 {
		fmt.Println(yellowStyle.Render("No offerings found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n📦 Offerings:\n"))
	for _, o := range offerings {
		status := greenStyle.Render("Current")
		if !o.IsCurrent {
			status = "Inactive"
		}
		fmt.Printf("  ID: %s\n", cyanStyle.Render(o.ID))
		fmt.Printf("  Identifier: %s\n", o.Identifier)
		fmt.Printf("  Name: %s\n", o.DisplayName)
		fmt.Printf("  Status: %s\n", status)
		fmt.Printf("  Packages: %d\n", len(o.Packages))
		fmt.Println()
	}
	fmt.Println(grayStyle.Render(fmt.Sprintf("Total: %d offerings", len(offerings))))

	return nil
}

func runInternalOfferingGet(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	offeringID, _ := cmd.Flags().GetString("offering-id")
	if offeringID == "" {
		return fmt.Errorf("offering-id is required (--offering-id or -o)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Fetching offering...")

	path := fmt.Sprintf("/developers/me/projects/%s/offerings/%s", projectID, offeringID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var offering internal.Offering
	if err := json.Unmarshal(toJSON(resp.Data), &offering); err != nil {
		return fmt.Errorf("error parsing offering: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📦 Offering Details:\n"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(offering.ID))
	fmt.Printf("  Identifier: %s\n", offering.Identifier)
	fmt.Printf("  Name: %s\n", offering.DisplayName)
	fmt.Printf("  Archived: %v\n", offering.IsArchived)
	fmt.Printf("  Current: %v\n", offering.IsCurrent)
	if offering.Metadata != nil {
		fmt.Printf("  Metadata: %v\n", offering.Metadata)
	}
	fmt.Printf("  Packages: %d\n", len(offering.Packages))

	if len(offering.Packages) > 0 {
		fmt.Println(internalStyle.Render("\n  Packages:"))
		for _, pkg := range offering.Packages {
			fmt.Printf("    - %s (%s)\n", cyanStyle.Render(pkg.Identifier), pkg.DisplayName)
		}
	}

	return nil
}

func runInternalOfferingsCreate(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
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

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Creating offering...")

	path := fmt.Sprintf("/developers/me/projects/%s/offerings", projectID)
	data := map[string]string{
		"identifier":   identifier,
		"display_name": name,
	}
	resp, err := client.Post(path, data)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	var offering internal.Offering
	if err := json.Unmarshal(toJSON(resp.Data), &offering); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(greenStyle.Render("\n✓ Offering created:"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(offering.ID))
	fmt.Printf("  Identifier: %s\n", offering.Identifier)
	fmt.Printf("  Name: %s\n", offering.DisplayName)

	return nil
}

func runInternalOfferingsDelete(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	offeringID, _ := cmd.Flags().GetString("offering-id")
	if offeringID == "" {
		return fmt.Errorf("offering-id is required (--offering-id or -o)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Deleting offering...")

	path := fmt.Sprintf("/developers/me/projects/%s/offerings/%s", projectID, offeringID)
	_, err = client.Delete(path)
	if err != nil {
		return err
	}

	fmt.Println(greenStyle.Render("\n✓ Offering deleted: " + offeringID))
	return nil
}

func runInternalOfferingsDuplicate(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	offeringID, _ := cmd.Flags().GetString("offering-id")
	identifier, _ := cmd.Flags().GetString("identifier")
	name, _ := cmd.Flags().GetString("name")
	packagesOnly, _ := cmd.Flags().GetBool("packages-only")

	if offeringID == "" {
		return fmt.Errorf("offering-id is required (--offering-id or -o)")
	}
	if identifier == "" {
		return fmt.Errorf("identifier is required (--identifier or -i)")
	}
	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Duplicating offering...")

	path := fmt.Sprintf("/developers/me/projects/%s/offerings/%s/duplicate", projectID, offeringID)
	data := map[string]interface{}{
		"identifier":     identifier,
		"display_name":  name,
		"packages_only": packagesOnly,
	}

	resp, err := client.Post(path, data)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	var offering map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &offering); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(greenStyle.Render("\n✓ Offering duplicated:"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(offering["id"].(string)))
	fmt.Printf("  Identifier: %s\n", offering["identifier"])
	fmt.Printf("  Name: %s\n", offering["display_name"])

	return nil
}

func runInternalOfferingsSetCurrent(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	offeringID, _ := cmd.Flags().GetString("offering-id")
	if offeringID == "" {
		return fmt.Errorf("offering-id is required (--offering-id or -o)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📦 Setting offering as current...")

	path := fmt.Sprintf("/developers/me/projects/%s/offerings/%s", projectID, offeringID)
	data := map[string]interface{}{
		"is_current": true,
	}

	resp, err := client.Patch(path, data)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	fmt.Println(greenStyle.Render("\n✓ Offering set as current"))

	return nil
}

func runInternalProductsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")

	fmt.Println("\n🛍️ Fetching products...")

	path := fmt.Sprintf("/developers/me/projects/%s/products", projectID)
	params := map[string]string{"limit": fmt.Sprintf("%d", limit)}
	resp, err := client.GetWithParams(path, params)
	if err != nil {
		return err
	}

	var products []internal.Product
	if err := json.Unmarshal(toJSON(resp.Items), &products); err != nil {
		return fmt.Errorf("error parsing products: %w", err)
	}

	if len(products) == 0 {
		fmt.Println(yellowStyle.Render("No products found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n🛍️ Products:\n"))
	for _, p := range products {
		fmt.Printf("  ID: %s\n", cyanStyle.Render(p.ID))
		fmt.Printf("  Identifier: %s\n", p.Identifier)
		fmt.Printf("  Type: %s\n", p.ProductType)
		if p.App != nil {
			fmt.Printf("  App: %s\n", p.App.Name)
		}
		fmt.Println()
	}
	fmt.Println(grayStyle.Render(fmt.Sprintf("Total: %d products", len(products))))

	return nil
}

func runInternalAppsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📱 Fetching apps...")

	path := fmt.Sprintf("/developers/me/projects/%s/apps", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var apps []internal.App
	if err := json.Unmarshal(toJSON(resp.Items), &apps); err != nil {
		return fmt.Errorf("error parsing apps: %w", err)
	}

	if len(apps) == 0 {
		fmt.Println(yellowStyle.Render("No apps found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n📱 Apps:\n"))
	for _, a := range apps {
		fmt.Printf("  ID: %s\n", cyanStyle.Render(a.ID))
		fmt.Printf("  Name: %s\n", a.Name)
		fmt.Printf("  Type: %s\n", a.Type)
		fmt.Println()
	}

	return nil
}

func runInternalProductStoresStatuses(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🏪 Fetching product stores statuses...")

	path := fmt.Sprintf("/developers/me/projects/%s/product_stores_statuses", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var statuses []map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Items), &statuses); err != nil {
		// Try parsing as Data if Items is empty
		var data map[string]interface{}
		if err2 := json.Unmarshal(toJSON(resp.Data), &data); err2 != nil {
			return fmt.Errorf("error parsing store statuses: %w", err)
		}
		fmt.Println(internalStyle.Render("\n🏪 Product Stores Status:\n"))
		for key, value := range data {
			fmt.Printf("  %s: %v\n", cyanStyle.Render(key), value)
		}
		return nil
	}

	if len(statuses) == 0 {
		fmt.Println(yellowStyle.Render("No store statuses found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n🏪 Product Stores Status:\n"))
	for _, s := range statuses {
		store := s["store"].(string)
		status := s["status"].(string)
		synced := s["last_synced_at"]
		fmt.Printf("  %s: %s", cyanStyle.Render(store), greenStyle.Render(status))
		if synced != nil {
			fmt.Printf(" (synced: %s)", synced)
		}
		fmt.Println()
	}

	return nil
}

func runInternalCollaboratorsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n👥 Fetching collaborators...")

	path := fmt.Sprintf("/developers/me/projects/%s/collaborators", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	collaborators, _ := data["collaborators"].([]interface{})
	owner, _ := data["owner"].(map[string]interface{})

	fmt.Println(internalStyle.Render("\n👥 Team:\n"))

	if owner != nil {
		fmt.Printf("  Owner: %s\n", cyanStyle.Render(owner["email"].(string)))
	}

	for _, c := range collaborators {
		collab := c.(map[string]interface{})
		fmt.Printf("  %s (%s)\n", collab["email"], collab["role"])
	}

	return nil
}

func runInternalAPIKeysList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔑 Fetching API keys...")

	path := fmt.Sprintf("/developers/me/projects/%s/api_keys", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var keys []map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Items), &keys); err != nil {
		return fmt.Errorf("error parsing keys: %w", err)
	}

	if len(keys) == 0 {
		fmt.Println(yellowStyle.Render("No API keys found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n🔑 API Keys:\n"))
	for _, k := range keys {
		keyType := k["key_type"].(string)
		keyStyle := yellowStyle
		if keyType == "secret" {
			keyStyle = cyanStyle
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
}

func runInternalAuditList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n📜 Fetching audit logs...")

	path := fmt.Sprintf("/developers/me/projects/%s/audit_logs", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	logs, _ := data["data"].([]interface{})

	if len(logs) == 0 {
		fmt.Println(yellowStyle.Render("No audit logs found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n📜 Recent Audit Logs:\n"))
	for _, l := range logs {
		log := l.(map[string]interface{})
		actionType := log["action_type"].(string)
		targetType := log["target_type"].(string)
		target := log["target_identifier"].(string)
		occurredAt := log["occurred_at"].(string)

		fmt.Printf("  %s %s %s %s\n",
			cyanStyle.Render(actionType),
			grayStyle.Render(targetType),
			yellowStyle.Render(target),
			grayStyle.Render(occurredAt[:10]))
	}

	return nil
}

func runInternalPriceExperimentsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔬 Fetching price experiments...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments", projectID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	experiments, _ := data["experiments"].([]interface{})

	if len(experiments) == 0 {
		fmt.Println(yellowStyle.Render("No price experiments found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n🔬 Price Experiments:\n"))
	for _, e := range experiments {
		exp := e.(map[string]interface{})
		id := exp["id"].(string)
		name := exp["name"].(string)
		status := exp["status"].(string)
		offeringID, _ := exp["offering_id"].(string)

		fmt.Printf("  ID: %s\n", cyanStyle.Render(id))
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Status: %s\n", greenStyle.Render(status))
		fmt.Printf("  Offering ID: %s\n", yellowStyle.Render(offeringID))
		fmt.Println()
	}

	return nil
}

func runInternalPriceExperimentGet(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔬 Fetching price experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s", projectID, experimentID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var exp map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &exp); err != nil {
		return fmt.Errorf("error parsing experiment: %w", err)
	}

	fmt.Println(internalStyle.Render("\n🔬 Price Experiment Details:\n"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(exp["id"].(string)))
	fmt.Printf("  Name: %s\n", exp["name"])
	fmt.Printf("  Status: %s\n", greenStyle.Render(exp["status"].(string)))
	fmt.Printf("  Offering ID: %s\n", yellowStyle.Render(exp["offering_id"].(string)))

	if variants, ok := exp["variants"].([]interface{}); ok && len(variants) > 0 {
		fmt.Println(internalStyle.Render("\n  Variants:"))
		for i, v := range variants {
			variant := v.(map[string]interface{})
			fmt.Printf("    %d. %s - %s\n", i+1, cyanStyle.Render(variant["id"].(string)), variant["name"])
		}
	}

	return nil
}

func runInternalPriceExperimentCreate(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	name, _ := cmd.Flags().GetString("name")
	offeringA, _ := cmd.Flags().GetString("offering-a")
	offeringB, _ := cmd.Flags().GetString("offering-b")
	enrollment, _ := cmd.Flags().GetInt("enrollment")
	expType, _ := cmd.Flags().GetString("type")
	primaryMetric, _ := cmd.Flags().GetString("primary-metric")
	notes, _ := cmd.Flags().GetString("notes")

	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}
	if offeringA == "" {
		return fmt.Errorf("offering-a is required")
	}
	if offeringB == "" {
		return fmt.Errorf("offering-b is required")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔬 Creating price experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments", projectID)
	data := map[string]interface{}{
		"display_name":         name,
		"offering_a_id":         offeringA,
		"offering_b_id":         offeringB,
		"enrollment_percentage": enrollment,
		"experiment_type":       expType,
		"primary_metric":        primaryMetric,
		"secondary_metrics":     []string{"Conversion to paying", "Trials started", "Active subscribers"},
		"notes":                 notes,
		"targeting_conditions":  []interface{}{},
		"placements":            nil,
	}

	resp, err := client.Post(path, data)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	var exp map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &exp); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(greenStyle.Render("\n✓ Price experiment created:"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(exp["id"].(string)))
	fmt.Printf("  Name: %s\n", exp["name"])
	fmt.Printf("  Status: %s\n", greenStyle.Render(exp["status"].(string)))

	return nil
}

func runInternalPriceExperimentPause(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔬 Pausing experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s/pause", projectID, experimentID)
	resp, err := client.Post(path, nil)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	fmt.Println(greenStyle.Render("\n✓ Experiment paused"))
	return nil
}

func runInternalPriceExperimentResume(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔬 Resuming experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s/resume", projectID, experimentID)
	resp, err := client.Post(path, nil)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	fmt.Println(greenStyle.Render("\n✓ Experiment resumed"))
	return nil
}

func runInternalPriceExperimentStop(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🔬 Stopping experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s/stop", projectID, experimentID)
	resp, err := client.Post(path, nil)
	if err != nil {
		return err
	}

	if resp.Code != "" {
		return fmt.Errorf("error: %s - %s", resp.Code, resp.Message)
	}

	fmt.Println(greenStyle.Render("\n✓ Experiment stopped"))
	return nil
}

func runInternalUtilitiesCountries(cmd *cobra.Command, args []string) error {
	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n🌍 Fetching countries...")

	path := "/utilities/countries"
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var countries []map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Items), &countries); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	if len(countries) == 0 {
		fmt.Println(yellowStyle.Render("No countries found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n🌍 Supported Countries:\n"))
	for _, c := range countries {
		code := c["code"].(string)
		name := c["name"].(string)
		fmt.Printf("  %s - %s\n", cyanStyle.Render(code), name)
	}

	return nil
}

func runInternalSubscriberListsList(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")

	fmt.Println("\n👥 Fetching subscriber lists...")

	path := fmt.Sprintf("/developers/me/projects/%s/subscriber_lists", projectID)
	params := map[string]string{"limit": fmt.Sprintf("%d", limit)}
	resp, err := client.GetWithParams(path, params)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	lists, _ := data["data"].([]interface{})

	if len(lists) == 0 {
		fmt.Println(yellowStyle.Render("No subscriber lists found."))
		return nil
	}

	fmt.Println(internalStyle.Render("\n👥 Subscriber Lists:\n"))
	for _, l := range lists {
		list := l.(map[string]interface{})
		id := list["id"].(string)
		name := list["name"].(string)
		size, _ := list["size"].(float64)
		listType, _ := list["type"].(string)

		fmt.Printf("  ID: %s\n", cyanStyle.Render(id))
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Type: %s\n", greenStyle.Render(listType))
		fmt.Printf("  Size: %.0f\n", size)
		fmt.Println()
	}

	return nil
}

func runInternalSubscriberListGet(cmd *cobra.Command, args []string) error {
	listID, _ := cmd.Flags().GetString("list-id")
	if listID == "" {
		return fmt.Errorf("list-id is required (--list-id or -l)")
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n👥 Fetching subscriber list...")

	path := fmt.Sprintf("/developers/me/subscriber_lists/%s", listID)
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var list map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &list); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n👥 Subscriber List Details:\n"))
	fmt.Printf("  ID: %s\n", cyanStyle.Render(list["id"].(string)))
	fmt.Printf("  Name: %s\n", list["name"])
	fmt.Printf("  Type: %s\n", greenStyle.Render(list["type"].(string)))
	if size, ok := list["size"].(float64); ok {
		fmt.Printf("  Size: %.0f\n", size)
	}

	return nil
}

func runInternalSubscriberListManifest(cmd *cobra.Command, args []string) error {
	client, err := getInternalClient()
	if err != nil {
		return err
	}

	fmt.Println("\n👥 Fetching subscriber lists manifest...")

	path := "/developers/me/subscriber_lists/manifest"
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n👥 Subscriber Lists Manifest:\n"))

	if totals, ok := data["totals"].(map[string]interface{}); ok {
		fmt.Println("  Totals:")
		for key, value := range totals {
			fmt.Printf("    %s: %v\n", cyanStyle.Render(key), value)
		}
	}

	if lists, ok := data["lists"].([]interface{}); ok && len(lists) > 0 {
		fmt.Println(internalStyle.Render("\n  Lists:"))
		for _, l := range lists {
			list := l.(map[string]interface{})
			fmt.Printf("    - %s (%s): %.0f\n",
				cyanStyle.Render(list["id"].(string)),
				list["name"],
				list["size"].(float64))
		}
	}

	return nil
}

func runInternalChartsOverview(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	fmt.Println("\n📊 Fetching project overview...")

	path := fmt.Sprintf("/developers/me/charts_v2/overview?app_uuid=%s&sandbox_mode=%t", projectID, sandbox)
	if appUUID != "" {
		path = fmt.Sprintf("/developers/me/charts_v2/overview?app_uuid=%s&sandbox_mode=%t", appUUID, sandbox)
	}

	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📊 Project Overview:\n"))

	if summary, ok := data["summary"].(map[string]interface{}); ok {
		fmt.Println("  Summary:")
		for key, value := range summary {
			fmt.Printf("    %s: %v\n", cyanStyle.Render(key), value)
		}
	}

	if charts, ok := data["charts"].([]interface{}); ok && len(charts) > 0 {
		fmt.Println(internalStyle.Render("\n  Charts Available:"))
		for _, c := range charts {
			chart := c.(map[string]interface{})
			fmt.Printf("    - %s (%s)\n",
				chart["title"],
				cyanStyle.Render(chart["type"].(string)))
		}
	}

	return nil
}

func runInternalChartsTrials(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	resolution, _ := cmd.Flags().GetInt("resolution")
	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	fmt.Println("\n📊 Fetching trial analytics...")

	path := fmt.Sprintf("/developers/me/charts_v2/trials?app_uuid=%s&sandbox_mode=%t&resolution=%d", projectID, sandbox, resolution)
	if appUUID != "" {
		path = fmt.Sprintf("/developers/me/charts_v2/trials?app_uuid=%s&sandbox_mode=%t&resolution=%d", appUUID, sandbox, resolution)
	}
	if startDate != "" {
		path += "&start_date=" + startDate
	}
	if endDate != "" {
		path += "&end_date=" + endDate
	}
	path += "&is_sparkline=true"

	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📊 Trial Analytics:\n"))

	if units, ok := data["units"].([]interface{}); ok && len(units) > 0 {
		fmt.Println("  Data Points:")
		for i, u := range units {
			unit := u.(map[string]interface{})
			if i < 10 { // Show first 10
				fmt.Printf("    %s: %v\n",
					cyanStyle.Render(unit["date"].(string)),
					unit["value"])
			}
		}
		if len(units) > 10 {
			fmt.Printf("    ... and %d more data points\n", len(units)-10)
		}
	}

	if total, ok := data["total"].(float64); ok {
		fmt.Printf("\n  Total Trials: %s\n", cyanStyle.Render(fmt.Sprintf("%.0f", total)))
	}

	return nil
}

func runInternalChartsOverviewAll(cmd *cobra.Command, args []string) error {
	client, err := getInternalClient()
	if err != nil {
		return err
	}

	sandbox, _ := cmd.Flags().GetBool("sandbox")

	fmt.Println("\n📊 Fetching overview for all projects...")

	path := fmt.Sprintf("/developers/me/charts_v2/overview?sandbox_mode=%t&v3=false", sandbox)

	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📊 All Projects Overview:\n"))

	if projects, ok := data["projects"].([]interface{}); ok && len(projects) > 0 {
		for _, p := range projects {
			project := p.(map[string]interface{})
			name := project["name"].(string)
			id := project["id"].(string)
			fmt.Printf("  Project: %s (%s)\n", cyanStyle.Render(name), cyanStyle.Render(id))
			if summary, ok := project["summary"].(map[string]interface{}); ok {
				for key, value := range summary {
					fmt.Printf("    %s: %v\n", key, value)
				}
			}
			fmt.Println()
		}
	}

	return nil
}

func runInternalChartsTransactions(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	resolution, _ := cmd.Flags().GetInt("resolution")
	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	fmt.Println("\n📊 Fetching transaction analytics...")

	path := fmt.Sprintf("/developers/me/charts_v2/transactions?app_uuid=%s&sandbox_mode=%t&resolution=%d", projectID, sandbox, resolution)
	if appUUID != "" {
		path = fmt.Sprintf("/developers/me/charts_v2/transactions?app_uuid=%s&sandbox_mode=%t&resolution=%d", appUUID, sandbox, resolution)
	}
	if startDate != "" {
		path += "&start_date=" + startDate
	}
	if endDate != "" {
		path += "&end_date=" + endDate
	}
	path += "&is_sparkline=true"

	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📊 Transaction Analytics:\n"))

	if units, ok := data["units"].([]interface{}); ok && len(units) > 0 {
		fmt.Println("  Data Points:")
		for i, u := range units {
			unit := u.(map[string]interface{})
			if i < 10 {
				fmt.Printf("    %s: %v\n",
					cyanStyle.Render(unit["date"].(string)),
					unit["value"])
			}
		}
		if len(units) > 10 {
			fmt.Printf("    ... and %d more data points\n", len(units)-10)
		}
	}

	if total, ok := data["total"].(float64); ok {
		fmt.Printf("\n  Total Transactions: %s\n", cyanStyle.Render(fmt.Sprintf("%.0f", total)))
	}

	return nil
}

func runInternalChartsRevenue(cmd *cobra.Command, args []string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}

	client, err := getInternalClient()
	if err != nil {
		return err
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	resolution, _ := cmd.Flags().GetInt("resolution")
	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	fmt.Println("\n📊 Fetching revenue analytics...")

	path := fmt.Sprintf("/developers/me/charts_v2/revenue?app_uuid=%s&sandbox_mode=%t&resolution=%d", projectID, sandbox, resolution)
	if appUUID != "" {
		path = fmt.Sprintf("/developers/me/charts_v2/revenue?app_uuid=%s&sandbox_mode=%t&resolution=%d", appUUID, sandbox, resolution)
	}
	if startDate != "" {
		path += "&start_date=" + startDate
	}
	if endDate != "" {
		path += "&end_date=" + endDate
	}
	path += "&is_sparkline=true"

	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(toJSON(resp.Data), &data); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(internalStyle.Render("\n📊 Revenue Analytics:\n"))

	if units, ok := data["units"].([]interface{}); ok && len(units) > 0 {
		fmt.Println("  Data Points:")
		for i, u := range units {
			unit := u.(map[string]interface{})
			if i < 10 {
				fmt.Printf("    %s: %v\n",
					cyanStyle.Render(unit["date"].(string)),
					unit["value"])
			}
		}
		if len(units) > 10 {
			fmt.Printf("    ... and %d more data points\n", len(units)-10)
		}
	}

	if total, ok := data["total"].(float64); ok {
		fmt.Printf("\n  Total Revenue: %s\n", cyanStyle.Render(fmt.Sprintf("%.2f", total)))
	}

	return nil
}

func toJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
