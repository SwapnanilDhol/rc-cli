package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var chartsCmd = &cobra.Command{
	Use:   "charts",
	Short: "View analytics charts (Internal Dashboard API)",
	Long:  `View analytics charts using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var chartsOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Get project overview analytics",
	RunE:  runChartsOverview,
}

var chartsOverviewAllCmd = &cobra.Command{
	Use:   "overview-all",
	Short: "Get overview analytics across all projects",
	RunE:  runChartsOverviewAll,
}

var chartsTrialsCmd = &cobra.Command{
	Use:   "trials",
	Short: "Get trial analytics",
	RunE:  runChartsTrials,
}

var chartsTransactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Get transaction analytics",
	RunE:  runChartsTransactions,
}

var chartsRevenueCmd = &cobra.Command{
	Use:   "revenue",
	Short: "Get revenue analytics",
	RunE:  runChartsRevenue,
}

func init() {
	chartsCmd.AddCommand(chartsOverviewCmd, chartsOverviewAllCmd, chartsTrialsCmd, chartsTransactionsCmd, chartsRevenueCmd)

	chartsOverviewCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsOverviewCmd.Flags().String("app-uuid", "", "App UUID (optional, uses first app if not specified)")
	chartsOverviewAllCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsTrialsCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	chartsTrialsCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	chartsTrialsCmd.Flags().Int("resolution", 0, "Resolution (0=daily, 1=weekly, 2=monthly)")
	chartsTrialsCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsTrialsCmd.Flags().String("app-uuid", "", "App UUID")
	chartsTransactionsCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	chartsTransactionsCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	chartsTransactionsCmd.Flags().Int("resolution", 0, "Resolution (0=daily, 1=weekly, 2=monthly)")
	chartsTransactionsCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsTransactionsCmd.Flags().String("app-uuid", "", "App UUID")
	chartsRevenueCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	chartsRevenueCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	chartsRevenueCmd.Flags().Int("resolution", 0, "Resolution (0=daily, 1=weekly, 2=monthly)")
	chartsRevenueCmd.Flags().Bool("sandbox", false, "Include sandbox data")
	chartsRevenueCmd.Flags().String("app-uuid", "", "App UUID")
}

func runChartsOverview(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	Progress("\n📊 Fetching project overview...")

	path := fmt.Sprintf("/developers/me/charts_v2/overview?app_uuid=%s&sandbox_mode=%t", projectID, sandbox)
	if appUUID != "" {
		path = fmt.Sprintf("/developers/me/charts_v2/overview?app_uuid=%s&sandbox_mode=%t", appUUID, sandbox)
	}

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

	fmt.Println(InternalStyle.Render("\n📊 Project Overview:\n"))

	if summary, ok := data["summary"].(map[string]interface{}); ok {
		fmt.Println("  Summary:")
		for key, value := range summary {
			fmt.Printf("    %s: %v\n", CyanStyle.Render(key), value)
		}
	}

	if charts, ok := data["charts"].([]interface{}); ok && len(charts) > 0 {
		fmt.Println(InternalStyle.Render("\n  Charts Available:"))
		for _, c := range charts {
			chart := c.(map[string]interface{})
			fmt.Printf("    - %s (%s)\n",
				chart["title"],
				CyanStyle.Render(chart["type"].(string)))
		}
	}

	return nil
}

func runChartsOverviewAll(cmd *cobra.Command, args []string) error {
	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	sandbox, _ := cmd.Flags().GetBool("sandbox")

	Progress("\n📊 Fetching overview for all projects...")

	path := fmt.Sprintf("/developers/me/charts_v2/overview?sandbox_mode=%t&v3=false", sandbox)

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

	fmt.Println(InternalStyle.Render("\n📊 All Projects Overview:\n"))

	if projects, ok := data["projects"].([]interface{}); ok && len(projects) > 0 {
		for _, p := range projects {
			project := p.(map[string]interface{})
			name := project["name"].(string)
			id := project["id"].(string)
			fmt.Printf("  Project: %s (%s)\n", CyanStyle.Render(name), CyanStyle.Render(id))
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

func runChartsTrials(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	resolution, _ := cmd.Flags().GetInt("resolution")
	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	Progress("\n📊 Fetching trial analytics...")

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

	fmt.Println(InternalStyle.Render("\n📊 Trial Analytics:\n"))

	if units, ok := data["units"].([]interface{}); ok && len(units) > 0 {
		fmt.Println("  Data Points:")
		for i, u := range units {
			unit := u.(map[string]interface{})
			if i < 10 {
				fmt.Printf("    %s: %v\n",
					CyanStyle.Render(unit["date"].(string)),
					unit["value"])
			}
		}
		if len(units) > 10 {
			fmt.Printf("    ... and %d more data points\n", len(units)-10)
		}
	}

	if total, ok := data["total"].(float64); ok {
		fmt.Printf("\n  Total Trials: %s\n", CyanStyle.Render(fmt.Sprintf("%.0f", total)))
	}

	return nil
}

func runChartsTransactions(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	resolution, _ := cmd.Flags().GetInt("resolution")
	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	Progress("\n📊 Fetching transaction analytics...")

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

	fmt.Println(InternalStyle.Render("\n📊 Transaction Analytics:\n"))

	if units, ok := data["units"].([]interface{}); ok && len(units) > 0 {
		fmt.Println("  Data Points:")
		for i, u := range units {
			unit := u.(map[string]interface{})
			if i < 10 {
				fmt.Printf("    %s: %v\n",
					CyanStyle.Render(unit["date"].(string)),
					unit["value"])
			}
		}
		if len(units) > 10 {
			fmt.Printf("    ... and %d more data points\n", len(units)-10)
		}
	}

	if total, ok := data["total"].(float64); ok {
		fmt.Printf("\n  Total Transactions: %s\n", CyanStyle.Render(fmt.Sprintf("%.0f", total)))
	}

	return nil
}

func runChartsRevenue(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	resolution, _ := cmd.Flags().GetInt("resolution")
	sandbox, _ := cmd.Flags().GetBool("sandbox")
	appUUID, _ := cmd.Flags().GetString("app-uuid")

	Progress("\n📊 Fetching revenue analytics...")

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

	fmt.Println(InternalStyle.Render("\n📊 Revenue Analytics:\n"))

	if units, ok := data["units"].([]interface{}); ok && len(units) > 0 {
		fmt.Println("  Data Points:")
		for i, u := range units {
			unit := u.(map[string]interface{})
			if i < 10 {
				fmt.Printf("    %s: %v\n",
					CyanStyle.Render(unit["date"].(string)),
					unit["value"])
			}
		}
		if len(units) > 10 {
			fmt.Printf("    ... and %d more data points\n", len(units)-10)
		}
	}

	if total, ok := data["total"].(float64); ok {
		fmt.Printf("\n  Total Revenue: %s\n", CyanStyle.Render(fmt.Sprintf("%.2f", total)))
	}

	return nil
}
