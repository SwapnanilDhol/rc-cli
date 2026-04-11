package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Email           string `mapstructure:"email"`
	Password        string `mapstructure:"password"`
	AuthToken       string `mapstructure:"authToken"`
	ProjectID       string `mapstructure:"projectId"`
	APIKey          string `mapstructure:"apiKey"` // Public v2 API key (legacy)
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}

	configPath := filepath.Join(homeDir, ".revenuerc")
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	viper.SetDefault("email", "")
	viper.SetDefault("password", "")
	viper.SetDefault("authToken", "")
	viper.SetDefault("projectId", "")
	viper.SetDefault("apiKey", "")

	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, nil
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}

	configDir := filepath.Dir(filepath.Join(homeDir, ".revenuerc"))
	os.MkdirAll(configDir, 0700)

	configPath := filepath.Join(homeDir, ".revenuerc")

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")
	viper.Set("email", cfg.Email)
	viper.Set("password", cfg.Password)
	viper.Set("authToken", cfg.AuthToken)
	viper.Set("projectId", cfg.ProjectID)
	viper.Set("apiKey", cfg.APIKey)

	return viper.WriteConfig()
}

func ClearAuth() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	cfg.AuthToken = ""
	cfg.Email = ""
	cfg.Password = ""
	return SaveConfig(cfg)
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/.revenuerc"
	}
	return filepath.Join(homeDir, ".revenuerc")
}
