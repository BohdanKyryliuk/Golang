package currency_converter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"Golang/currencyapi"

	"github.com/joho/godotenv"
)

// client is the CurrencyAPI client instance
var client *currencyapi.Client

// CurrencyConverterError wraps errors from the currency converter
type CurrencyConverterError struct {
	Operation string
	Err       error
}

func (e *CurrencyConverterError) Error() string {
	return e.Operation + ": " + e.Err.Error()
}

func (e *CurrencyConverterError) Unwrap() error {
	return e.Err
}

// initCurrencyApi initializes the currency API client
func initCurrencyApi() error {
	if client != nil {
		return nil // Already initialized
	}

	// Load .env file
	loadEnv()

	// Get API key from environment variable
	apiKey := os.Getenv("CURRENCY_API_KEY")
	if apiKey == "" {
		return &CurrencyConverterError{
			Operation: "init",
			Err:       errors.New("CURRENCY_API_KEY not set in .env file or environment variables"),
		}
	}

	var err error
	client, err = currencyapi.NewClient(
		apiKey,
		currencyapi.WithTimeout(15*time.Second),
	)
	if err != nil {
		return &CurrencyConverterError{
			Operation: "create_client",
			Err:       err,
		}
	}

	return nil
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, will use environment variables")
	}
}

// CheckStatus returns the API status or an error
func CheckStatus() (string, error) {
	if err := initCurrencyApi(); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := client.Status(ctx)
	if err != nil {
		return "", handleAPIError("check_status", err)
	}

	jsonBytes, err := json.Marshal(status)
	if err != nil {
		return "", &CurrencyConverterError{
			Operation: "marshal_status",
			Err:       err,
		}
	}
	return string(jsonBytes), nil
}

// GetCurrencies returns available currencies or an error
func GetCurrencies() (string, error) {
	if err := initCurrencyApi(); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currencies, err := client.Currencies(ctx, nil)
	if err != nil {
		return "", handleAPIError("get_currencies", err)
	}

	jsonBytes, err := json.Marshal(currencies)
	if err != nil {
		return "", &CurrencyConverterError{
			Operation: "marshal_currencies",
			Err:       err,
		}
	}
	return string(jsonBytes), nil
}

// GetLatestRates returns the latest exchange rates or an error
func GetLatestRates() (string, error) {
	if err := initCurrencyApi(); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	latestRates, err := client.Latest(ctx, &currencyapi.LatestParams{
		BaseCurrency: "USD",
		Currencies:   []string{"UAH", "EUR"},
	})
	if err != nil {
		return "", handleAPIError("get_latest_rates", err)
	}

	jsonBytes, err := json.Marshal(latestRates)
	if err != nil {
		return "", &CurrencyConverterError{
			Operation: "marshal_latest_rates",
			Err:       err,
		}
	}
	return string(jsonBytes), nil
}

// handleAPIError processes API errors and wraps them appropriately
func handleAPIError(operation string, err error) error {
	// Log detailed error information
	var apiErr *currencyapi.APIError
	if errors.As(err, &apiErr) {
		log.Printf("[%s] API error: code=%s, message=%s, status=%d",
			operation, apiErr.Code, apiErr.Message, apiErr.StatusCode)
	}

	var httpErr *currencyapi.HTTPError
	if errors.As(err, &httpErr) {
		log.Printf("[%s] HTTP error: status=%d, body=%s",
			operation, httpErr.StatusCode, httpErr.Body)
	}

	var reqErr *currencyapi.RequestError
	if errors.As(err, &reqErr) {
		log.Printf("[%s] Request error during %s: %v",
			operation, reqErr.Op, reqErr.Err)
	}

	return &CurrencyConverterError{
		Operation: operation,
		Err:       err,
	}
}
