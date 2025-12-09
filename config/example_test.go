// Package config_test provides examples of using the config package
package config_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"Golang/config"
	"Golang/currency_converter"
)

// Example_directConfig shows how to create a currency converter client
// with direct configuration (not from environment variables)
func Example_directConfig() {
	// Create the client with explicit configuration
	client, err := currency_converter.New(currency_converter.Config{
		APIKey:  "your-api-key-here",
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use the client
	_ = client
	fmt.Println("Client created successfully")
	// Output: Client created successfully
}

// Example_fromEnvironment shows how to create a currency converter client
// loading configuration from environment variables
func Example_fromEnvironment() {
	// Set environment variable for demonstration
	os.Setenv("CURRENCY_API_KEY", "demo-api-key")

	// Create the client from environment
	client, err := currency_converter.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use the client
	_ = client
	fmt.Println("Client created from environment")
	// Output: Client created from environment
}

// Example_customConfigLoader shows how to use the config package directly
// for more control over configuration loading
func Example_customConfigLoader() {
	// Set environment variable for demonstration
	os.Setenv("CURRENCY_API_KEY", "custom-api-key")

	// Load configuration using the config package
	cfg, err := config.LoadCurrencyAPIConfig(nil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// You can modify the config before creating the client
	modifiedTimeout := cfg.Timeout * 2

	// Create the client with modified config
	client, err := currency_converter.New(currency_converter.Config{
		APIKey:  cfg.APIKey,
		Timeout: modifiedTimeout,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	_ = client
	fmt.Println("Client created with custom config")
	// Output: Client created with custom config
}
