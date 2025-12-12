package worker

import (
	"context"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.Currencies) == 0 {
		t.Error("DefaultConfig should have default currencies")
	}

	if cfg.FetchInterval == 0 {
		t.Error("DefaultConfig should have a default fetch interval")
	}

	if cfg.RequestTimeout == 0 {
		t.Error("DefaultConfig should have a default request timeout")
	}
}

func TestNewManagerRequiresClient(t *testing.T) {
	_, err := NewManager(nil, DefaultConfig())
	if err == nil {
		t.Error("NewManager should return error when apiClient is nil")
	}
}

func TestRateStore(t *testing.T) {
	store := NewRateStore()

	// Test Get on empty store
	_, err := store.Get("USD")
	if err == nil {
		t.Error("Get should return error for non-existent currency")
	}

	var notFoundErr *NotFoundError
	if !errorAs(err, &notFoundErr) {
		t.Error("Expected NotFoundError")
	}

	// Test Set and Get
	rateData := &RateData{
		BaseCurrency: "USD",
		FetchedAt:    time.Now(),
	}
	store.Set("USD", rateData)

	got, err := store.Get("USD")
	if err != nil {
		t.Errorf("Get returned error: %v", err)
	}

	if got.BaseCurrency != "USD" {
		t.Errorf("Expected USD, got %s", got.BaseCurrency)
	}

	// Test GetAll
	all := store.GetAll()
	if len(all) != 1 {
		t.Errorf("Expected 1 item, got %d", len(all))
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Currency: "XYZ"}
	expected := "rates not found for currency: XYZ"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestManagerStartStop(t *testing.T) {
	// This is an integration test that would require a mock API client
	// For now, we just test the basic flow
	t.Skip("Skipping integration test - requires mock API client")
}

// Helper function for error type checking
func errorAs(err error, target interface{}) bool {
	if err == nil {
		return false
	}
	switch v := target.(type) {
	case **NotFoundError:
		e, ok := err.(*NotFoundError)
		if ok {
			*v = e
			return true
		}
	}
	return false
}

func TestConfigDefaults(t *testing.T) {
	// Test that empty config gets populated with defaults
	cfg := Config{}

	// When creating a manager, defaults should be applied
	// We can't test this directly without an API client, but we can
	// verify DefaultConfig returns valid values

	defaults := DefaultConfig()

	if defaults.FetchInterval != 1*time.Minute {
		t.Errorf("Expected 1 minute fetch interval, got %v", defaults.FetchInterval)
	}

	if defaults.RequestTimeout != 10*time.Second {
		t.Errorf("Expected 10 second request timeout, got %v", defaults.RequestTimeout)
	}

	// Check that config has nil values before defaults
	if cfg.FetchInterval != 0 {
		t.Error("Empty config should have zero FetchInterval")
	}
}

func TestRateStoreThreadSafety(t *testing.T) {
	store := NewRateStore()

	// Run concurrent reads and writes
	done := make(chan bool)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				done <- true
				return
			default:
				store.Set("USD", &RateData{
					BaseCurrency: "USD",
					FetchedAt:    time.Now(),
				})
			}
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				done <- true
				return
			default:
				_, _ = store.Get("USD")
				_ = store.GetAll()
			}
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done
}
