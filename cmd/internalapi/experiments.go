package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var experimentsCmd = &cobra.Command{
	Use:   "experiments",
	Short: "Manage price experiments (A/B tests)",
	Long:  `Manage price experiments using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var experimentsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List price experiments",
	RunE:    runPriceExperimentsList,
}

var experimentsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get price experiment details",
	RunE:  runPriceExperimentGet,
}

var experimentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a price experiment",
	RunE:  runPriceExperimentCreate,
}

var experimentsPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause a running experiment",
	RunE:  runPriceExperimentPause,
}

var experimentsResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused experiment",
	RunE:  runPriceExperimentResume,
}

var experimentsStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop an experiment",
	RunE:  runPriceExperimentStop,
}

var experimentsTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List available price experiment types",
	RunE:  runExperimentTypes,
}

var experimentTypes = map[string]string{
	"introductory_offer":    "Test introductory offers (free trial, pay as you go)",
	"free_trial_offer":      "Test free trial variations",
	"paywall_design":        "Test different paywall layouts and designs",
	"price_point":           "Test different price points for the same product",
	"subscription_duration": "Test different subscription durations",
	"subscription_ordering": "Test different package orderings",
	"other":                 "Custom experiment type",
}

func init() {
	experimentsCmd.AddCommand(experimentsListCmd, experimentsGetCmd, experimentsCreateCmd, experimentsPauseCmd, experimentsResumeCmd, experimentsStopCmd, experimentsTypesCmd)

	experimentsGetCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID")

	experimentsCreateCmd.Flags().StringP("name", "n", "", "Experiment display name (required)")
	experimentsCreateCmd.Flags().StringP("offering-a", "a", "", "Offering A ID (required)")
	experimentsCreateCmd.Flags().StringP("offering-b", "b", "", "Offering B ID (required)")
	experimentsCreateCmd.Flags().Int("enrollment", 100, "Enrollment percentage (default 100)")
	experimentsCreateCmd.Flags().String("type", "paywall_design", "Experiment type (see: rc internal experiments types)")
	experimentsCreateCmd.Flags().String("primary-metric", "Realized LTV per customer", "Primary metric")
	experimentsCreateCmd.Flags().String("notes", "", "Experiment notes")

	experimentsPauseCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID (required)")
	experimentsResumeCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID (required)")
	experimentsStopCmd.Flags().StringP("experiment-id", "e", "", "Experiment ID (required)")
}

func runExperimentTypes(cmd *cobra.Command, args []string) error {
	fmt.Println(InternalStyle.Render("\n📋 Available Price Experiment Types:\n"))
	for t, desc := range experimentTypes {
		fmt.Printf("  %s  %s\n", CyanStyle.Render(t), desc)
	}
	fmt.Println()
	return nil
}

func runPriceExperimentsList(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n🔬 Fetching price experiments...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments", projectID)
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

	experiments := resp.Items

	if len(experiments) == 0 {
		fmt.Println(YellowStyle.Render("No price experiments found."))
		return nil
	}

	fmt.Println(InternalStyle.Render("\n🔬 Price Experiments:\n"))
	for _, e := range experiments {
		exp := e.(map[string]interface{})
		id, _ := exp["id"].(string)
		name, _ := exp["display_name"].(string)
		running, _ := exp["is_running"].(bool)
		expType, _ := exp["experiment_type"].(string)
		startingAt, _ := exp["starting_at"].(string)
		createdAt, _ := exp["created_at"].(string)

		fmt.Printf("  ID: %s\n", CyanStyle.Render(id))
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Type: %s\n", expType)
		fmt.Printf("  Running: %v\n", running)
		fmt.Printf("  Starting: %s\n", startingAt)
		fmt.Printf("  Created: %s\n", createdAt)
		fmt.Println()
	}

	return nil
}

func runPriceExperimentGet(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n🔬 Fetching price experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s", projectID, experimentID)
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

	var exp map[string]interface{}
	if err := json.Unmarshal(ToJSON(resp.Data), &exp); err != nil {
		return fmt.Errorf("error parsing experiment: %w", err)
	}

	fmt.Println(InternalStyle.Render("\n🔬 Price Experiment Details:\n"))
	fmt.Printf("  ID: %s\n", CyanStyle.Render(exp["id"].(string)))
	fmt.Printf("  Name: %s\n", exp["display_name"])
	fmt.Printf("  Type: %s\n", exp["experiment_type"])
	fmt.Printf("  Running: %v\n", exp["is_running"])
	fmt.Printf("  Enrollment: %v%%\n", exp["enrollment_percentage"])

	if offeringA, ok := exp["offering_a"].(map[string]interface{}); ok {
		fmt.Printf("  Offering A: %s (%s)\n", YellowStyle.Render(offeringA["display_name"].(string)), CyanStyle.Render(offeringA["id"].(string)))
	}
	if offeringB, ok := exp["offering_b"].(map[string]interface{}); ok {
		fmt.Printf("  Offering B: %s (%s)\n", YellowStyle.Render(offeringB["display_name"].(string)), CyanStyle.Render(offeringB["id"].(string)))
	}

	if variants, ok := exp["variants"].([]interface{}); ok && len(variants) > 0 {
		fmt.Println(InternalStyle.Render("\n  Variants:"))
		for i, v := range variants {
			variant := v.(map[string]interface{})
			fmt.Printf("    %d. %s - %s\n", i+1, CyanStyle.Render(variant["id"].(string)), variant["name"])
		}
	}

	return nil
}

func runPriceExperimentCreate(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
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

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n🔬 Creating price experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments", projectID)
	data := map[string]interface{}{
		"display_name":          name,
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
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	var exp map[string]interface{}
	if err := json.Unmarshal(ToJSON(resp.Data), &exp); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println(GreenStyle.Render("\n✓ Price experiment created:"))
	fmt.Printf("  ID: %s\n", CyanStyle.Render(exp["id"].(string)))
	fmt.Printf("  Name: %s\n", exp["display_name"])
	fmt.Printf("  Running: %v\n", exp["is_running"])

	return nil
}

func runPriceExperimentPause(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n🔬 Pausing experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s/pause", projectID, experimentID)
	resp, err := client.Post(path, nil)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Experiment paused"))
	return nil
}

func runPriceExperimentResume(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n🔬 Resuming experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s/resume", projectID, experimentID)
	resp, err := client.Post(path, nil)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Experiment resumed"))
	return nil
}

func runPriceExperimentStop(cmd *cobra.Command, args []string) error {
	projectID, err := GetProjectID()
	if err != nil {
		return err
	}

	experimentID, _ := cmd.Flags().GetString("experiment-id")
	if experimentID == "" {
		return fmt.Errorf("experiment-id is required (--experiment-id or -e)")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	Progress("\n🔬 Stopping experiment...")

	path := fmt.Sprintf("/developers/me/projects/%s/price_experiments/%s/stop", projectID, experimentID)
	resp, err := client.Post(path, nil)
	if err != nil {
		return err
	}
	if err := CheckResponse(resp); err != nil {
		return err
	}
	if EmitJSON(resp) {
		return nil
	}

	fmt.Println(GreenStyle.Render("\n✓ Experiment stopped"))
	return nil
}
