// Package worker provides a currency rate fetching system with configurable workers.
// Each worker is responsible for fetching rates for a specific base currency.
package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/BohdanKyryliuk/golang/currencyapi"
)

// RateData holds the cached rate information for a currency
type RateData struct {
	BaseCurrency  string                          `json:"base_currency"`
	Rates         map[string]currencyapi.RateInfo `json:"rates"`
	LastUpdatedAt string                          `json:"last_updated_at"`
	FetchedAt     time.Time                       `json:"fetched_at"`
}

// Config holds the configuration for the worker manager
type Config struct {
	// Currencies is the list of base currencies to fetch rates for (one worker per currency)
	Currencies []string
	// FetchInterval is the interval between rate fetches (default: 1 minute)
	FetchInterval time.Duration
	// RequestTimeout is the timeout for individual API requests (default: 10 seconds)
	RequestTimeout time.Duration
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		Currencies:     []string{"USD", "EUR", "GBP"},
		FetchInterval:  1 * time.Minute,
		RequestTimeout: 10 * time.Second,
	}
}

// Manager manages multiple currency rate workers
type Manager struct {
	config    Config
	apiClient *currencyapi.Client
	store     *RateStore
	workers   []*Worker
	stopCh    chan struct{}
	wg        sync.WaitGroup
	running   bool
	mu        sync.RWMutex
}

// NewManager creates a new worker manager
func NewManager(apiClient *currencyapi.Client, cfg Config) (*Manager, error) {
	if apiClient == nil {
		return nil, errors.New("API client is required")
	}

	defaults := DefaultConfig()
	if len(cfg.Currencies) == 0 {
		cfg.Currencies = defaults.Currencies
	}
	if cfg.FetchInterval == 0 {
		cfg.FetchInterval = defaults.FetchInterval
	}
	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = defaults.RequestTimeout
	}
	return &Manager{
		config:    cfg,
		apiClient: apiClient,
		store:     NewRateStore(),
		stopCh:    make(chan struct{}),
	}, nil
}

// Start starts all workers
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return errors.New("workers are already running")
	}

	log.Printf("Starting %d currency rate workers", len(m.config.Currencies))

	for _, currency := range m.config.Currencies {
		worker := NewWorker(currency, m.apiClient, m.store, m.config)
		m.workers = append(m.workers, worker)

		m.wg.Add(1)
		go func(w *Worker) {
			defer m.wg.Done()
			w.Run(ctx, m.stopCh)
		}(worker)
	}

	m.running = true
	log.Println("All workers started")
	return nil
}

// Stop stops all workers gracefully
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	log.Println("Stopping all workers...")
	close(m.stopCh)
	m.wg.Wait()
	m.running = false
	log.Println("All workers stopped")
}

// GetRates returns the cached rates for a specific base currency
func (m *Manager) GetRates(baseCurrency string) (*RateData, error) {
	return m.store.Get(baseCurrency)
}

// GetAllRates returns all cached rates
func (m *Manager) GetAllRates() map[string]*RateData {
	return m.store.GetAll()
}

// GetCurrencies returns the list of currencies being tracked
func (m *Manager) GetCurrencies() []string {
	return m.config.Currencies
}

// IsRunning returns whether the workers are running
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// Worker fetches rates for a specific base currency
type Worker struct {
	baseCurrency string
	apiClient    *currencyapi.Client
	store        *RateStore
	config       Config
}

// NewWorker creates a new worker for a specific currency
func NewWorker(baseCurrency string, apiClient *currencyapi.Client, store *RateStore, cfg Config) *Worker {
	return &Worker{
		baseCurrency: baseCurrency,
		apiClient:    apiClient,
		store:        store,
		config:       cfg,
	}
}

// Run starts the worker's fetch loop
func (w *Worker) Run(ctx context.Context, stopCh <-chan struct{}) {
	log.Printf("[%s] Worker started", w.baseCurrency)

	// Fetch immediately on start
	w.fetch(ctx)

	ticker := time.NewTicker(w.config.FetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Worker stopped: context cancelled", w.baseCurrency)
			return
		case <-stopCh:
			log.Printf("[%s] Worker stopped: stop signal received", w.baseCurrency)
			return
		case <-ticker.C:
			w.fetch(ctx)
		}
	}
}

// fetch fetches the latest rates and stores them
func (w *Worker) fetch(ctx context.Context) {
	fetchCtx, cancel := context.WithTimeout(ctx, w.config.RequestTimeout)
	defer cancel()

	log.Printf("[%s] Fetching latest rates...", w.baseCurrency)

	response, err := w.apiClient.Latest(fetchCtx, &currencyapi.LatestParams{
		BaseCurrency: w.baseCurrency,
	})
	if err != nil {
		log.Printf("[%s] Error fetching rates: %v", w.baseCurrency, err)
		return
	}

	rateData := &RateData{
		BaseCurrency:  w.baseCurrency,
		Rates:         response.Data,
		LastUpdatedAt: response.Meta.LastUpdatedAt,
		FetchedAt:     time.Now(),
	}

	w.store.Set(w.baseCurrency, rateData)
	log.Printf("[%s] Updated rates: %d currencies", w.baseCurrency, len(response.Data))
}

// RateStore is a thread-safe storage for currency rates
type RateStore struct {
	data map[string]*RateData
	mu   sync.RWMutex
}

// NewRateStore creates a new rate store
func NewRateStore() *RateStore {
	return &RateStore{
		data: make(map[string]*RateData),
	}
}

// Set stores rate data for a currency
func (s *RateStore) Set(currency string, data *RateData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[currency] = data
}

// Get retrieves rate data for a currency
func (s *RateStore) Get(currency string) (*RateData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.data[currency]
	if !ok {
		return nil, &NotFoundError{Currency: currency}
	}
	return data, nil
}

// GetAll returns all stored rate data
func (s *RateStore) GetAll() map[string]*RateData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]*RateData, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

// NotFoundError is returned when rate data is not found for a currency
type NotFoundError struct {
	Currency string
}

func (e *NotFoundError) Error() string {
	return "rates not found for currency: " + e.Currency
}
