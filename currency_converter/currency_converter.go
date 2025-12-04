package currency_converter

import (
	"log"
	"os"

	"github.com/everapihq/currencyapi-go"
	"github.com/joho/godotenv"
)

func initCurrencyApi() {
	// Load .env file
	loadEnv()

	// Get API key from environment variable
	apiKey := os.Getenv("CURRENCY_API_KEY")
	if apiKey == "" {
		log.Fatal("CURRENCY_API_KEY not set in .env file or environment variables")
	}

	currencyapi.Init(apiKey)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, will use environment variables")
	}
}

func CheckStatus() string {
	initCurrencyApi()
	status := currencyapi.Status()
	return string(status)
}

func GetCurrencies() string {
	initCurrencyApi()
	currencies := currencyapi.Currencies(map[string]string{
		"format": "json",
	})
	return string(currencies)
}

func GetLatestRates() string {
	initCurrencyApi()
	latestRates := currencyapi.Latest(map[string]string{
		"base_currency": "USD",
		"currencies":    "UAH,EUR",
	})
	return string(latestRates)
}
