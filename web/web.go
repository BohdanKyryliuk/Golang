package web

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BohdanKyryliuk/golang/currency_converter"
	"github.com/BohdanKyryliuk/golang/http/handler"
	"github.com/BohdanKyryliuk/golang/worker"
	"github.com/gin-gonic/gin"
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

	// Create Gin router
	router := gin.Default()

	// Register basic routes
	router.GET("/", handler.Hello)
	router.GET("/count", handler.Counter)
	router.POST("/count", handler.Counter)

	var workerManager *worker.Manager

	// Only register currency handlers if the client is available
	if cfg.CurrencyClient != nil {
		currencyHandler := handler.NewCurrency(cfg.CurrencyClient)

		// Create currency route group
		currencyGroup := router.Group("/currency")
		{
			currencyGroup.GET("/status", currencyHandler.Status)
			currencyGroup.GET("/currencies", currencyHandler.Currencies)
			currencyGroup.GET("/latest", currencyHandler.LatestRates)
		}

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

					// Create rates route group
					ratesGroup := router.Group("/rates")
					{
						ratesGroup.GET("", ratesHandler.GetRate)
						ratesGroup.GET("/all", ratesHandler.GetAllRates)
						ratesGroup.GET("/status", ratesHandler.GetWorkerStatus)
					}
				}
			}
		}
	}

	// Create HTTP server
	server := &gin.Engine{}
	*server = *router

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

		cancel()
	}()

	log.Println("Listening on :3001")

	if err := router.Run(":3001"); err != nil {
		log.Fatal(err)
	}
}
