package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initCharts() {
	chartsCmd := &cobra.Command{
		Use:   "charts",
		Short: "View analytics and charts",
	}

	chartsRevenueCmd := &cobra.Command{
		Use:   "revenue",
		Short: "Get revenue analytics",
		RunE:  runRevenueChart,
	}
	chartsRevenueCmd.Flags().StringP("start-date", "s", "", "Start date (YYYY-MM-DD)")
	chartsRevenueCmd.Flags().StringP("end-date", "e", "", "End date (YYYY-MM-DD)")

	chartsSubscribersCmd := &cobra.Command{
		Use:   "subscribers",
		Short: "Get subscriber analytics",
		RunE:  runSubscribersChart,
	}
	chartsSubscribersCmd.Flags().StringP("start-date", "s", "", "Start date (YYYY-MM-DD)")
	chartsSubscribersCmd.Flags().StringP("end-date", "e", "", "End date (YYYY-MM-DD)")

	chartsMRRCmd := &cobra.Command{
		Use:   "mrr",
		Short: "Get Monthly Recurring Revenue (MRR) analytics",
		RunE:  runMRRChart,
	}

	chartsChurnCmd := &cobra.Command{
		Use:   "churn",
		Short: "Get churn analytics",
		RunE:  runChurnChart,
	}

	chartsSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Get project summary/overview",
		RunE:  runSummary,
	}

	chartsCmd.AddCommand(chartsRevenueCmd, chartsSubscribersCmd, chartsMRRCmd, chartsChurnCmd, chartsSummaryCmd)
	RootCmd.AddCommand(chartsCmd)

	// Aliases
	RootCmd.AddCommand(&cobra.Command{
		Use:   "analytics:revenue",
		Short: "Get revenue analytics",
		RunE:  runRevenueChart,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use:   "analytics:mrr",
		Short: "Get MRR analytics",
		RunE:  runMRRChart,
	})
}

func runSummary(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n📊 Fetching project summary...")

	appsResp, err := client.Get(fmt.Sprintf("/projects/%s/apps", cfg.ProjectID))
	if err != nil {
		return err
	}

	customersResp, err := client.GetWithParams(fmt.Sprintf("/projects/%s/customers", cfg.ProjectID), map[string]string{"limit": "1"})
	if err != nil {
		return err
	}

	offeringsResp, err := client.Get(fmt.Sprintf("/projects/%s/offerings", cfg.ProjectID))
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n📊 Project Summary\n"))
	fmt.Println(cyanStyle.Render(fmt.Sprintf("  Apps: %d", len(appsResp.Items))))
	fmt.Println(cyanStyle.Render("  Customers: (See dashboard for full count)"))
	fmt.Println(cyanStyle.Render(fmt.Sprintf("  Offerings: %d", len(offeringsResp.Items))))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("  For detailed analytics, visit the RevenueCat dashboard:"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("  https://app.revenuecat.com/%s\n", cfg.ProjectID)))

	_ = customersResp

	return nil
}

func runRevenueChart(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	fmt.Println(appsStyle.Render("\n💰 Revenue Analytics\n"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Detailed revenue data requires RevenueCat Pro plan"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("  Visit: https://app.revenuecat.com/%s/revenue\n", cfg.ProjectID)))

	return nil
}

func runSubscribersChart(cmd *cobra.Command, args []string) error {
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

	fmt.Println("\n👥 Fetching subscriber data...")

	resp, err := client.GetWithParams(fmt.Sprintf("/projects/%s/customers", cfg.ProjectID), map[string]string{"limit": "100"})
	if err != nil {
		return err
	}

	fmt.Println(appsStyle.Render("\n👥 Subscriber Analytics\n"))

	platforms := make(map[string]int)
	for _, item := range resp.Items {
		sub := item.(map[string]interface{})
		platform := getStringValue(sub, "last_seen_platform", "Unknown")
		platforms[platform]++
	}

	fmt.Println(cyanStyle.Render("  By Platform:"))
	for platform, count := range platforms {
		fmt.Printf("    %s: %d\n", platform, count)
	}
	fmt.Println()

	return nil
}

func runMRRChart(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	fmt.Println(appsStyle.Render("\n📈 Monthly Recurring Revenue (MRR)\n"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: MRR analytics available in RevenueCat dashboard"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("  Visit: https://app.revenuecat.com/%s/mrr\n", cfg.ProjectID)))

	return nil
}

func runChurnChart(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ProjectID == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("\n⚠ No project ID configured. Run: rc config"))
		return nil
	}

	fmt.Println(appsStyle.Render("\n📉 Churn Analytics\n"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("  Note: Churn analytics available in RevenueCat dashboard"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(fmt.Sprintf("  Visit: https://app.revenuecat.com/%s/churn\n", cfg.ProjectID)))

	return nil
}
