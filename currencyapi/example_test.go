package currencyapi_test

import (
	"Golang/currencyapi"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// ExampleNewClient demonstrates how to create a new CurrencyAPI client
func ExampleNewClient() {
	// Create a basic client with API key
	client, err := currencyapi.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// Or create a client with custom options
	client, err = currencyapi.NewClient(
		"your-api-key",
		currencyapi.WithTimeout(30*time.Second),
		currencyapi.WithBaseURL("https://api.currencyapi.com/v3/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = client // use client
}

// ExampleClient_Latest demonstrates fetching latest exchange rates with error handling
func ExampleClient_Latest() {
	apiKey := os.Getenv("CURRENCY_API_KEY")
	client, err := currencyapi.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch latest rates
	response, err := client.Latest(ctx, &currencyapi.LatestParams{
		BaseCurrency: "USD",
		Currencies:   []string{"EUR", "UAH", "GBP"},
	})

	if err != nil {
		// Handle different error types
		handleError(err)
		return
	}

	// Use the response
	for code, rate := range response.Data {
		fmt.Printf("%s: %f\n", code, rate.Value)
	}
}

// ExampleClient_Convert demonstrates currency conversion with error handling
func ExampleClient_Convert() {
	apiKey := os.Getenv("CURRENCY_API_KEY")
	client, err := currencyapi.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Convert 100 USD to EUR and UAH
	response, err := client.Convert(ctx, &currencyapi.ConvertParams{
		BaseCurrency: "USD",
		Currencies:   []string{"EUR", "UAH"},
		Value:        100.0,
	})

	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("Converting 100 USD:\n")
	for code, rate := range response.Data {
		fmt.Printf("  %s: %f\n", code, rate.Value)
	}
}

// Example_errorHandling demonstrates comprehensive error handling following Go best practices
func Example_errorHandling() {
	apiKey := os.Getenv("CURRENCY_API_KEY")
	client, err := currencyapi.NewClient(apiKey)
	if err != nil {
		// Handle client creation error (e.g., missing API key)
		var validationErr *currencyapi.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf("Configuration error: %s - %s\n", validationErr.Field, validationErr.Message)
		}
		return
	}

	ctx := context.Background()
	_, err = client.Latest(ctx, nil)

	if err != nil {
		// Method 1: Use type assertion with errors.As (recommended)
		var apiErr *currencyapi.APIError
		if errors.As(err, &apiErr) {
			// Handle API-specific errors
			if apiErr.IsInvalidAPIKey() {
				fmt.Println("Please check your API key")
			} else if apiErr.IsQuotaExceeded() {
				fmt.Println("Monthly quota exceeded, please upgrade your plan")
			} else {
				fmt.Printf("API error: %s\n", apiErr.Message)
			}
			return
		}

		var httpErr *currencyapi.HTTPError
		if errors.As(err, &httpErr) {
			// Handle HTTP-level errors
			if httpErr.IsNotFound() {
				fmt.Println("Endpoint not found")
			} else if httpErr.IsRateLimited() {
				fmt.Println("Rate limited, please wait and retry")
			} else {
				fmt.Printf("HTTP error %d: %s\n", httpErr.StatusCode, httpErr.Body)
			}
			return
		}

		var reqErr *currencyapi.RequestError
		if errors.As(err, &reqErr) {
			// Handle request-level errors (network issues, etc.)
			fmt.Printf("Request failed during %s: %v\n", reqErr.Op, reqErr.Err)

			// Check for context cancellation
			if errors.Is(reqErr.Err, context.Canceled) {
				fmt.Println("Request was cancelled")
			} else if errors.Is(reqErr.Err, context.DeadlineExceeded) {
				fmt.Println("Request timed out")
			}
			return
		}

		var parseErr *currencyapi.ParseError
		if errors.As(err, &parseErr) {
			fmt.Printf("Failed to parse response from %s: %v\n", parseErr.Endpoint, parseErr.Err)
			return
		}

		// Fallback for unknown errors
		fmt.Printf("Unknown error: %v\n", err)
	}
}

// Example_retryLogic demonstrates implementing retry logic for temporary errors
func Example_retryLogic() {
	apiKey := os.Getenv("CURRENCY_API_KEY")
	client, err := currencyapi.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	var response *currencyapi.LatestResponse
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err = client.Latest(ctx, &currencyapi.LatestParams{
			BaseCurrency: "USD",
		})

		if err == nil {
			break // Success
		}

		// Check if error is temporary and worth retrying
		if currencyapi.IsTemporaryError(err) {
			waitTime := time.Duration(attempt+1) * time.Second
			fmt.Printf("Temporary error (attempt %d/%d), retrying in %v: %v\n",
				attempt+1, maxRetries, waitTime, err)
			time.Sleep(waitTime)
			continue
		}

		// Non-temporary error, don't retry
		break
	}

	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("Successfully fetched %d rates\n", len(response.Data))
}

// handleError is a helper function that handles different error types
func handleError(err error) {
	// Method 2: Use helper functions for simpler checks
	if currencyapi.IsValidationError(err) {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	if currencyapi.IsAPIError(err) {
		fmt.Printf("API error: %v\n", err)
		// Get status code if available
		if statusCode, ok := currencyapi.GetHTTPStatusCode(err); ok {
			fmt.Printf("Status code: %d\n", statusCode)
		}
		return
	}

	if currencyapi.IsHTTPError(err) {
		fmt.Printf("HTTP error: %v\n", err)
		return
	}

	if currencyapi.IsRequestError(err) {
		fmt.Printf("Request error: %v\n", err)
		return
	}

	if currencyapi.IsParseError(err) {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Printf("Unknown error: %v\n", err)
}
