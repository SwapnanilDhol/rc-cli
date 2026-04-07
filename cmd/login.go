package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/config"
	"revenuecat-cli/internal"
)

func init() {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to RevenueCat dashboard API (session-based auth)",
		Long:  `Authenticate with email and password to get a session cookie for the internal dashboard API.

Usage:
  rc internal login --email user@example.com --password secret

Or run interactively:
  rc internal login`,
		RunE:  runLogin,
	}
	loginCmd.Flags().String("email", "", "Email address")
	loginCmd.Flags().String("password", "", "Password")
	internalCmd.AddCommand(loginCmd)

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from RevenueCat dashboard API",
		RunE:  runLogout,
	}
	internalCmd.AddCommand(logoutCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")

	// If email/password not provided via flags, use interactive mode
	if email == "" || password == "" {
		return runLoginInteractive()
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	fmt.Println("\n🔐 Logging in...")

	loginResp, err := internal.Login(email, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	cfg.Email = email
	cfg.Password = password
	cfg.AuthToken = loginResp.AuthenticationToken

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("\n✓ Logged in as " + loginResp.Email))
	return nil
}

func runLoginInteractive() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Check if already logged in
	if cfg.AuthToken != "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Already logged in as " + cfg.Email))
		return nil
	}

	qs := []*survey.Question{
		{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Email:",
				Default: cfg.Email,
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); ok && len(str) == 0 {
					return fmt.Errorf("email is required")
				}
				return nil
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password:",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); ok && len(str) == 0 {
					return fmt.Errorf("password is required")
				}
				return nil
			},
		},
	}

	answers := struct {
		Email    string
		Password string
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	fmt.Println("\n🔐 Logging in...")

	loginResp, err := internal.Login(answers.Email, answers.Password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	cfg.Email = answers.Email
	cfg.Password = answers.Password
	cfg.AuthToken = loginResp.AuthenticationToken

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("\n✓ Logged in as " + loginResp.Email))
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := config.ClearAuth(); err != nil {
		return fmt.Errorf("error clearing auth: %w", err)
	}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("\n✓ Logged out"))
	return nil
}
