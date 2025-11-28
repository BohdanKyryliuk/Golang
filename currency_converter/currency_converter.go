package currency_converter

import (
	"github.com/everapihq/currencyapi-go"
)

func initCurrencyApi() {
	currencyapi.Init("cur_live_4izJseiR0m05pAtCFprAMJ06aROnF5pEqmHpccLE")
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
