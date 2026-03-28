package cmd

import (
	"fmt"
	"os"
	"time"

	"revenuecat-cli/config"
)

func loadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Command-line flags override config file
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	if projectID != "" {
		cfg.ProjectID = projectID
	}

	return cfg, nil
}

func getStringValue(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return fallback
}

func getBoolValue(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func formatDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return "N/A"
}

func formatDateTime(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return "N/A"
}

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return home
}

func getConfigPath() string {
	return fmt.Sprintf("%s/.revenuerc", userHomeDir())
}

func formatUnixTime(ts int64) string {
	if ts == 0 {
		return "N/A"
	}
	return time.Unix(ts/1000, 0).Format("2006-01-02")
}
