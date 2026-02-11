package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ProviderConfig struct {
	Name      string `mapstructure:"name"`       // Display name (e.g., "Groq", "My Local LLM")
	Type      string `mapstructure:"type"`       // "openai", "anthropic", "gemini"
	Model     string `mapstructure:"model"`      // e.g., "gpt-4o", "claude-3-opus"
	APIKeyEnv string `mapstructure:"api_key_env"`// Env var name (e.g., "OPENAI_API_KEY")
	BaseURL   string `mapstructure:"base_url"`   // Optional: for OpenAI-compatible APIs
}

type Config struct {
	Providers []ProviderConfig `mapstructure:"providers"`
}

// Default chain used if no config file is found
var defaultProviders = []ProviderConfig{
	{Name: "Groq", Type: "openai", Model: "llama-3.3-70b-versatile", BaseURL: "https://api.groq.com/openai/v1", APIKeyEnv: "GROQ_API_KEY"},
	{Name: "Gemini", Type: "gemini", Model: "gemini-2.5-flash", APIKeyEnv: "GEMINI_API_KEY"},
	{Name: "Anthropic", Type: "anthropic", Model: "claude-3-5-sonnet-20240620", APIKeyEnv: "ANTHROPIC_API_KEY"},
	{Name: "OpenAI", Type: "openai", Model: "gpt-4o-mini", BaseURL: "https://api.openai.com/v1", APIKeyEnv: "OPENAI_API_KEY"},
}

func Load() (*Config, error) {
	v := viper.New()
	// Will find config.yaml, config.toml, config.json, etc (Viper supports several file formats)
	v.SetConfigName("config")
	
	// 1. Check local directory
	v.AddConfigPath(".")
	
	// 2. Check standard OS config directory
	configDir, err := os.UserConfigDir()
	if err == nil {
		v.AddConfigPath(filepath.Join(configDir, "error-explain"))
	}

	// 3. Try to read
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No config file found? Use defaults silently.
			return &Config{Providers: defaultProviders}, nil
		}
		// Config file exists but has errors
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// If the config file exists but is empty/missing providers, fallback to default
	if len(cfg.Providers) == 0 {
		return &Config{Providers: defaultProviders}, nil
	}

	return &cfg, nil
}
