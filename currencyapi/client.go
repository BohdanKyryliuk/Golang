// Package currencyapi provides a robust client for the CurrencyAPI service
// with proper error handling following Go best practices.
// See: https://go.dev/blog/error-handling-and-go
package currencyapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default API endpoint
	DefaultBaseURL = "https://api.currencyapi.com/v3/"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 10 * time.Second
)

// Client represents a CurrencyAPI client with configurable options
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new CurrencyAPI client with the provided API key and options
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, &ValidationError{Field: "apiKey", Message: "API key is required"}
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// doRequest performs an HTTP request and returns the response body or an error
func (c *Client) doRequest(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	// Build URL with query parameters
	reqURL, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, &RequestError{
			Op:  "parse_url",
			Err: err,
		}
	}

	// Add query parameters
	q := reqURL.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	reqURL.RawQuery = q.Encode()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, &RequestError{
			Op:  "create_request",
			Err: err,
		}
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{
			Op:  "execute_request",
			Err: err,
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{
			Op:  "read_response",
			Err: err,
		}
	}

	// Check for API errors based on status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, c.parseAPIError(resp.StatusCode, body)
	}

	return body, nil
}

// parseAPIError attempts to parse an API error response
func (c *Client) parseAPIError(statusCode int, body []byte) error {
	var apiErr APIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If we can't parse the error, return a generic HTTP error
		return &HTTPError{
			StatusCode: statusCode,
			Body:       string(body),
		}
	}

	return &APIError{
		StatusCode: statusCode,
		Code:       apiErr.Error.Code,
		Message:    apiErr.Error.Message,
		Info:       apiErr.Error.Info,
	}
}

// Status returns the current API status
func (c *Client) Status(ctx context.Context) (*StatusResponse, error) {
	body, err := c.doRequest(ctx, "status", nil)
	if err != nil {
		return nil, err
	}

	var response StatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &ParseError{
			Endpoint: "status",
			Err:      err,
		}
	}

	return &response, nil
}

// Currencies returns available currencies
func (c *Client) Currencies(ctx context.Context, params *CurrenciesParams) (*CurrenciesResponse, error) {
	queryParams := make(map[string]string)
	if params != nil {
		if len(params.Currencies) > 0 {
			queryParams["currencies"] = strings.Join(params.Currencies, ",")
		}
		if params.Type != "" {
			queryParams["type"] = params.Type
		}
	}

	body, err := c.doRequest(ctx, "currencies", queryParams)
	if err != nil {
		return nil, err
	}

	var response CurrenciesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &ParseError{
			Endpoint: "currencies",
			Err:      err,
		}
	}

	return &response, nil
}

// Latest returns the latest exchange rates
func (c *Client) Latest(ctx context.Context, params *LatestParams) (*LatestResponse, error) {
	queryParams := make(map[string]string)
	if params != nil {
		if params.BaseCurrency != "" {
			queryParams["base_currency"] = params.BaseCurrency
		}
		if len(params.Currencies) > 0 {
			queryParams["currencies"] = strings.Join(params.Currencies, ",")
		}
	}

	body, err := c.doRequest(ctx, "latest", queryParams)
	if err != nil {
		return nil, err
	}

	var response LatestResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &ParseError{
			Endpoint: "latest",
			Err:      err,
		}
	}

	return &response, nil
}

// Historical returns historical exchange rates for a specific date
func (c *Client) Historical(ctx context.Context, params *HistoricalParams) (*HistoricalResponse, error) {
	if params == nil || params.Date == "" {
		return nil, &ValidationError{
			Field:   "date",
			Message: "date parameter is required for historical endpoint",
		}
	}

	queryParams := map[string]string{
		"date": params.Date,
	}
	if params.BaseCurrency != "" {
		queryParams["base_currency"] = params.BaseCurrency
	}
	if len(params.Currencies) > 0 {
		queryParams["currencies"] = strings.Join(params.Currencies, ",")
	}

	body, err := c.doRequest(ctx, "historical", queryParams)
	if err != nil {
		return nil, err
	}

	var response HistoricalResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &ParseError{
			Endpoint: "historical",
			Err:      err,
		}
	}

	return &response, nil
}

// Convert converts an amount from one currency to another
func (c *Client) Convert(ctx context.Context, params *ConvertParams) (*ConvertResponse, error) {
	if params == nil {
		return nil, &ValidationError{
			Field:   "params",
			Message: "convert parameters are required",
		}
	}

	queryParams := make(map[string]string)
	if params.BaseCurrency != "" {
		queryParams["base_currency"] = params.BaseCurrency
	}
	if len(params.Currencies) > 0 {
		queryParams["currencies"] = strings.Join(params.Currencies, ",")
	}
	if params.Value != 0 {
		queryParams["value"] = fmt.Sprintf("%f", params.Value)
	}
	if params.Date != "" {
		queryParams["date"] = params.Date
	}

	body, err := c.doRequest(ctx, "convert", queryParams)
	if err != nil {
		return nil, err
	}

	var response ConvertResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &ParseError{
			Endpoint: "convert",
			Err:      err,
		}
	}

	return &response, nil
}
