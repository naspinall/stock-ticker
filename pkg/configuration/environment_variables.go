package configuration

import (
	"fmt"
	"os"
	"strconv"
)

// Parses environment variables into a StockTickerConfiguration
type EnvironmentVariableProvider struct{}

func NewEnvironmentVariableProvider() *EnvironmentVariableProvider {
	return &EnvironmentVariableProvider{}
}

func (evp *EnvironmentVariableProvider) ParseConfiguration() (*StockTickerConfiguration, error) {

	apiKey, set := os.LookupEnv("STOCK_TICKER_API_KEY")
	if !set {
		return nil, fmt.Errorf("STOCK_TICKER_API_KEY not set")
	}

	numberOfDaysString, set := os.LookupEnv("STOCK_TICKER_DEFAULT_DAY_COUNT")
	if !set {
		return nil, fmt.Errorf("STOCK_TICKER_DEFAULT_DAY_COUNT not set")
	}

	symbol, set := os.LookupEnv("STOCK_TICKER_DEFAULT_SYMBOL")
	if !set {
		return nil, fmt.Errorf("STOCK_TICKER_DEFAULT_SYMBOL not set")
	}

	numberOfDays, err := strconv.Atoi(numberOfDaysString)
	if err != nil {
		return nil, fmt.Errorf("STOCK_TICKER_DEFAULT_DAY_COUNT not a valid integer, %w", err)
	}

	return &StockTickerConfiguration{
		AlphaVantageAPIKey:  apiKey,
		DefaultNumberOfDays: numberOfDays,
		DefaultSymbol:       symbol,
	}, nil
}
