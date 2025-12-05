package currencyapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		opts    []ClientOption
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "test-api-key",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: true,
		},
		{
			name:   "with custom timeout",
			apiKey: "test-api-key",
			opts: []ClientOption{
				WithTimeout(30 * time.Second),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClient_Latest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify API key header
		if r.Header.Get("apikey") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "invalid_api_key",
					"message": "Invalid API key",
				},
			})
			return
		}

		// Return mock response
		response := LatestResponse{
			Meta: struct {
				LastUpdatedAt string `json:"last_updated_at"`
			}{
				LastUpdatedAt: "2025-12-05T12:00:00Z",
			},
			Data: map[string]RateInfo{
				"EUR": {Code: "EUR", Value: 0.85},
				"UAH": {Code: "UAH", Value: 41.5},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", WithBaseURL(server.URL+"/"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	response, err := client.Latest(ctx, &LatestParams{
		BaseCurrency: "USD",
		Currencies:   []string{"EUR", "UAH"},
	})

	if err != nil {
		t.Fatalf("Latest() error = %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 currencies, got %d", len(response.Data))
	}

	if response.Data["EUR"].Value != 0.85 {
		t.Errorf("Expected EUR value 0.85, got %f", response.Data["EUR"].Value)
	}
}

func TestClient_APIError(t *testing.T) {
	// Create a mock server that returns an API error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "invalid_api_key",
				"message": "Your API key is invalid",
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient("invalid-key", WithBaseURL(server.URL+"/"))

	ctx := context.Background()
	_, err := client.Latest(ctx, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("Expected APIError, got %T", err)
	}

	if apiErr.Code != "invalid_api_key" {
		t.Errorf("Expected error code 'invalid_api_key', got '%s'", apiErr.Code)
	}

	if !apiErr.IsInvalidAPIKey() {
		t.Error("Expected IsInvalidAPIKey() to return true")
	}
}

func TestClient_RateLimitError(t *testing.T) {
	// Create a mock server that returns rate limit error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "quota_exceeded",
				"message": "You have exceeded your monthly quota",
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL+"/"))

	ctx := context.Background()
	_, err := client.Latest(ctx, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("Expected APIError, got %T", err)
	}

	if !apiErr.IsQuotaExceeded() {
		t.Error("Expected IsQuotaExceeded() to return true")
	}

	if !IsTemporaryError(err) {
		t.Error("Expected IsTemporaryError() to return true for rate limit error")
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	// Create a mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL+"/"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.Latest(ctx, nil)

	if err == nil {
		t.Fatal("Expected error due to context cancellation")
	}

	var reqErr *RequestError
	if !errors.As(err, &reqErr) {
		t.Errorf("Expected RequestError, got %T", err)
	}
}

func TestValidationError(t *testing.T) {
	client, _ := NewClient("test-key")

	ctx := context.Background()
	_, err := client.Historical(ctx, nil)

	if err == nil {
		t.Fatal("Expected validation error for missing date parameter")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	if validationErr.Field != "date" {
		t.Errorf("Expected validation error on 'date' field, got '%s'", validationErr.Field)
	}
}

func TestErrorHelperFunctions(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		isValidation   bool
		isRequest      bool
		isHTTP         bool
		isAPI          bool
		isParse        bool
		isTemporary    bool
		expectedStatus int
		hasStatus      bool
	}{
		{
			name:         "ValidationError",
			err:          &ValidationError{Field: "test", Message: "test"},
			isValidation: true,
		},
		{
			name:      "RequestError",
			err:       &RequestError{Op: "test", Err: errors.New("test")},
			isRequest: true,
		},
		{
			name:           "HTTPError 404",
			err:            &HTTPError{StatusCode: 404, Body: "not found"},
			isHTTP:         true,
			expectedStatus: 404,
			hasStatus:      true,
		},
		{
			name:           "HTTPError 500",
			err:            &HTTPError{StatusCode: 500, Body: "server error"},
			isHTTP:         true,
			isTemporary:    true,
			expectedStatus: 500,
			hasStatus:      true,
		},
		{
			name:           "APIError",
			err:            &APIError{StatusCode: 401, Code: "invalid_api_key", Message: "Invalid"},
			isAPI:          true,
			expectedStatus: 401,
			hasStatus:      true,
		},
		{
			name:    "ParseError",
			err:     &ParseError{Endpoint: "test", Err: errors.New("json error")},
			isParse: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsValidationError(tt.err) != tt.isValidation {
				t.Errorf("IsValidationError() = %v, want %v", IsValidationError(tt.err), tt.isValidation)
			}
			if IsRequestError(tt.err) != tt.isRequest {
				t.Errorf("IsRequestError() = %v, want %v", IsRequestError(tt.err), tt.isRequest)
			}
			if IsHTTPError(tt.err) != tt.isHTTP {
				t.Errorf("IsHTTPError() = %v, want %v", IsHTTPError(tt.err), tt.isHTTP)
			}
			if IsAPIError(tt.err) != tt.isAPI {
				t.Errorf("IsAPIError() = %v, want %v", IsAPIError(tt.err), tt.isAPI)
			}
			if IsParseError(tt.err) != tt.isParse {
				t.Errorf("IsParseError() = %v, want %v", IsParseError(tt.err), tt.isParse)
			}
			if IsTemporaryError(tt.err) != tt.isTemporary {
				t.Errorf("IsTemporaryError() = %v, want %v", IsTemporaryError(tt.err), tt.isTemporary)
			}

			status, hasStatus := GetHTTPStatusCode(tt.err)
			if hasStatus != tt.hasStatus {
				t.Errorf("GetHTTPStatusCode() hasStatus = %v, want %v", hasStatus, tt.hasStatus)
			}
			if hasStatus && status != tt.expectedStatus {
				t.Errorf("GetHTTPStatusCode() status = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHTTPErrorMethods(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		isNotFound     bool
		isUnauthorized bool
		isRateLimited  bool
	}{
		{"404 Not Found", 404, true, false, false},
		{"401 Unauthorized", 401, false, true, false},
		{"429 Too Many Requests", 429, false, false, true},
		{"500 Server Error", 500, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &HTTPError{StatusCode: tt.statusCode}

			if err.IsNotFound() != tt.isNotFound {
				t.Errorf("IsNotFound() = %v, want %v", err.IsNotFound(), tt.isNotFound)
			}
			if err.IsUnauthorized() != tt.isUnauthorized {
				t.Errorf("IsUnauthorized() = %v, want %v", err.IsUnauthorized(), tt.isUnauthorized)
			}
			if err.IsRateLimited() != tt.isRateLimited {
				t.Errorf("IsRateLimited() = %v, want %v", err.IsRateLimited(), tt.isRateLimited)
			}
		})
	}
}
