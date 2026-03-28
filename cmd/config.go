package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"revenuecat-cli/config"
)

func initConfig() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage RevenueCat CLI configuration",
		RunE:  runConfig,
	}
	RootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	qs := []*survey.Question{
		{
			Name: "apiKey",
			Prompt: &survey.Input{
				Message: "Enter your RevenueCat API key:",
				Default: cfg.APIKey,
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); ok && len(str) == 0 {
					return fmt.Errorf("API key is required")
				}
				return nil
			},
		},
		{
			Name: "projectId",
			Prompt: &survey.Input{
				Message: "Enter your RevenueCat Project ID (optional):",
				Default: cfg.ProjectID,
			},
		},
	}

	answers := struct {
		APIKey    string
		ProjectID string
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	cfg.APIKey = answers.APIKey
	cfg.ProjectID = answers.ProjectID
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Println("\n✓ Configuration saved successfully!")
	return nil
}
