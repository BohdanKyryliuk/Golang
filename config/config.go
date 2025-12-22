// Package config provides configuration loading utilities for the application.
// It separates config reading logic from business logic following best practices.
package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// CurrencyAPIConfig holds configuration for the CurrencyAPI client
type CurrencyAPIConfig struct {
	APIKey  string
	Timeout time.Duration
}

// ConfigError represents a configuration-related error
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error for " + e.Field + ": " + e.Message
}

// LoadOptions configures how configuration is loaded
type LoadOptions struct {
	EnvFile string // Path to .env file (optional)
}

// LoadCurrencyAPIConfig loads CurrencyAPI configuration from environment variables
func LoadCurrencyAPIConfig(opts *LoadOptions) (*CurrencyAPIConfig, error) {
	// Load .env file if specified or try default
	if opts != nil && opts.EnvFile != "" {
		if err := godotenv.Load(opts.EnvFile); err != nil {
			return nil, &ConfigError{
				Field:   "env_file",
				Message: "failed to load env file: " + err.Error(),
			}
		}
	} else {
		// Try to load default .env file, ignore if not found
		_ = godotenv.Load()
	}

	apiKey := os.Getenv("CURRENCY_API_KEY")
	if apiKey == "" {
		return nil, &ConfigError{
			Field:   "CURRENCY_API_KEY",
			Message: "environment variable is required but not set",
		}
	}

	// Default timeout
	timeout := 15 * time.Second

	return &CurrencyAPIConfig{
		APIKey:  apiKey,
		Timeout: timeout,
	}, nil
}

// Validate checks if the configuration is valid
func (c *CurrencyAPIConfig) Validate() error {
	if c.APIKey == "" {
		return errors.New("API key is required")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	return nil
}
