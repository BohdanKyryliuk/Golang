package currency_converter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/BohdanKyryliuk/golang/config"
	"github.com/BohdanKyryliuk/golang/currencyapi"
)

// Client represents a currency converter client with its dependencies
type Client struct {
	config    Config
	apiClient currencyapi.Client
}

// Config holds configuration for the currency converter
type Config struct {
	APIKey         string
	Timeout        time.Duration // HTTP client timeout
	RequestTimeout time.Duration // Individual request timeout (default: 10s)
}

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

// New creates a new Client from a Config struct
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, &CurrencyConverterError{
			Operation: "init",
			Err:       errors.New("API key is required"),
		}
	}

	// Set default timeout if not provided
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}

	// Set default request timeout if not provided
	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = 10 * time.Second
	}

	apiClient, err := currencyapi.NewHttpApiClient(
		cfg.APIKey,
		currencyapi.WithTimeout(cfg.Timeout),
	)
	if err != nil {
		return nil, &CurrencyConverterError{
			Operation: "create_client",
			Err:       err,
		}
	}

	return &Client{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

// APIClient returns the underlying currencyapi.Client for direct access
func (c *Client) APIClient() currencyapi.Client {
	return c.apiClient
}

// NewFromEnv creates a new Client by loading configuration from environment variables
// This is a convenience function that loads config from the config package
func NewFromEnv() (*Client, error) {
	cfg, err := config.LoadCurrencyAPIConfig(nil)
	if err != nil {
		return nil, &CurrencyConverterError{
			Operation: "load_config",
			Err:       err,
		}
	}

	return New(Config{
		APIKey:  cfg.APIKey,
		Timeout: cfg.Timeout,
	})
}

// CheckStatus returns the API status or an error
func (c *Client) CheckStatus(ctx context.Context) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.config.RequestTimeout)
		defer cancel()
	}

	status, err := c.apiClient.Status(ctx)
	if err != nil {
		return "", c.handleAPIError("check_status", err)
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
func (c *Client) GetCurrencies(ctx context.Context) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.config.RequestTimeout)
		defer cancel()
	}

	currencies, err := c.apiClient.Currencies(ctx, nil)
	if err != nil {
		return "", c.handleAPIError("get_currencies", err)
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

// LatestRatesParams holds parameters for fetching latest rates
type LatestRatesParams struct {
	BaseCurrency string   // Base currency code (default: USD)
	Currencies   []string // Target currency codes to fetch
}

// GetLatestRates returns the latest exchange rates or an error
func (c *Client) GetLatestRates(ctx context.Context, params *LatestRatesParams) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.config.RequestTimeout)
		defer cancel()
	}

	// Build API params with defaults
	apiParams := &currencyapi.LatestParams{}
	if params != nil {
		apiParams.BaseCurrency = params.BaseCurrency
		apiParams.Currencies = params.Currencies
	}
	if apiParams.BaseCurrency == "" {
		apiParams.BaseCurrency = "USD"
	}

	latestRates, err := c.apiClient.Latest(ctx, apiParams)
	if err != nil {
		return "", c.handleAPIError("get_latest_rates", err)
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
func (c *Client) handleAPIError(operation string, err error) error {
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
