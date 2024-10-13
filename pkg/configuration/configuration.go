package configuration

type StockTickerConfiguration struct {
	AlphaVantageAPIKey  string
	DefaultNumberOfDays int
	DefaultSymbol       string
}

// Generic provider to allow for reading from a file, another service etc.
type ConfigurationProvider interface {
	ParseConfiguration() (*StockTickerConfiguration, error)
}
