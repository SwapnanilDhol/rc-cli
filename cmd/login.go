package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"revenuecat-cli/config"
	"revenuecat-cli/internal"
)

func initAuth(root *cobra.Command) {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to the RevenueCat dashboard (enables 'rc internal' commands)",
		Long: `Authenticate with your RevenueCat email and password to obtain a dashboard
session. This is what 'rc internal …' commands use; the public v2 API commands
use an API key from 'rc config' instead.

Credentials and session token are stored in ~/.revenuerc.

Non-interactive:
  rc login --email you@example.com --password '…'
  RC_EMAIL=you@example.com RC_PASSWORD='…' rc login`,
		RunE: runLogin,
	}
	loginCmd.Flags().String("email", "", "RevenueCat email (or set RC_EMAIL)")
	loginCmd.Flags().String("password", "", "RevenueCat password (or set RC_PASSWORD)")
	loginCmd.Flags().Bool("force", false, "Re-authenticate even if a session already exists")
	root.AddCommand(loginCmd)

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear the stored dashboard session",
		RunE:  runLogout,
	}
	root.AddCommand(logoutCmd)
}

// envCredential reads RC_EMAIL/RC_PASSWORD, falling back to the older
// RC_INTERNAL_* names used before rc-internal was merged into rc.
func envCredential(name string) string {
	if v := os.Getenv("RC_" + name); v != "" {
		return v
	}
	return os.Getenv("RC_INTERNAL_" + name)
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")
	force, _ := cmd.Flags().GetBool("force")

	if email == "" {
		email = envCredential("EMAIL")
	}
	if password == "" {
		password = envCredential("PASSWORD")
	}

	// Explicit credentials always mean "log in as this account" — otherwise there
	// is no way to switch accounts without running logout first.
	if email != "" && password != "" {
		return performLogin(email, password)
	}

	if cfg.AuthToken != "" && !force {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).
			Render("Already logged in as " + cfg.Email + " (use --force to re-authenticate, or 'rc logout')"))
		return nil
	}

	if email == "" {
		email = cfg.Email
	}

	qs := []*survey.Question{
		{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Email:",
				Default: email,
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

	return performLogin(answers.Email, answers.Password)
}

func performLogin(email, password string) error {
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

func runLogout(cmd *cobra.Command, args []string) error {
	if err := config.ClearAuth(); err != nil {
		return fmt.Errorf("error clearing auth: %w", err)
	}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("\n✓ Logged out"))
	return nil
}
