package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BohdanKyryliuk/golang/currency_converter"
	"github.com/BohdanKyryliuk/golang/http/handler"
	"github.com/BohdanKyryliuk/golang/worker"
)

// ServerConfig holds configuration for the web server
type ServerConfig struct {
	// CurrencyClient is the currency converter client
	CurrencyClient *currency_converter.Client
	// WorkerConfig is the configuration for currency rate workers
	WorkerConfig *worker.Config
}

func StartServer() {
	// Initialize currency converter client from environment variables
	currencyClient, err := currency_converter.NewFromEnv()
	if err != nil {
		log.Printf("Warning: Currency converter not available: %v", err)
		// Continue without currency endpoints
	}

	// Default worker config
	var workerConfig *worker.Config
	if currencyClient != nil {
		cfg := worker.DefaultConfig()
		workerConfig = &cfg
	}

	StartServerWithConfig(ServerConfig{
		CurrencyClient: currencyClient,
		WorkerConfig:   workerConfig,
	})
}

// StartServerWithConfig starts the server with the provided configuration
func StartServerWithConfig(cfg ServerConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Hello)
	mux.HandleFunc("/count", handler.Counter)

	var workerManager *worker.Manager

	// Only register currency handlers if the client is available
	if cfg.CurrencyClient != nil {
		currencyHandler := handler.NewCurrency(cfg.CurrencyClient)
		mux.HandleFunc("/currency/status", currencyHandler.Status)
		mux.HandleFunc("/currency/currencies", currencyHandler.Currencies)
		mux.HandleFunc("/currency/latest", currencyHandler.LatestRates)

		// Initialize and start workers if config is provided
		if cfg.WorkerConfig != nil {
			var err error
			workerManager, err = worker.NewManager(cfg.CurrencyClient.APIClient(), *cfg.WorkerConfig)
			if err != nil {
				log.Printf("Warning: Failed to create worker manager: %v", err)
			} else {
				if err := workerManager.Start(ctx); err != nil {
					log.Printf("Warning: Failed to start workers: %v", err)
				} else {
					// Register rate handlers
					ratesHandler := handler.NewRates(workerManager)
					mux.HandleFunc("/rates", ratesHandler.GetRate)
					mux.HandleFunc("/rates/all", ratesHandler.GetAllRates)
					mux.HandleFunc("/rates/status", ratesHandler.GetWorkerStatus)
				}
			}
		}
	}

	server := &http.Server{
		Addr:    ":3001",
		Handler: mux,
	}

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down server...")

		// Stop workers first
		if workerManager != nil {
			workerManager.Stop()
		}

		// Shutdown HTTP server with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		cancel()
	}()

	log.Println("Listening on :3001")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
