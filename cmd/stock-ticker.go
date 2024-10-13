package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/naspinall/stock-ticker/pkg/configuration"
	"github.com/naspinall/stock-ticker/pkg/handlers"
	"github.com/naspinall/stock-ticker/pkg/middleware"
	"github.com/naspinall/stock-ticker/pkg/models"
)

func main() {

	// Parse configuration using a configuration provider
	var configurationProvider configuration.ConfigurationProvider = configuration.NewEnvironmentVariableProvider()
	tickerConfiguration, err := configurationProvider.ParseConfiguration()
	if err != nil {
		slog.Error("cannot parse configuration", "error", err)
		os.Exit(1)
	}

	// Initialize contextual logger for request tracing
	contextHandler := middleware.NewContextHandler(slog.NewJSONHandler(os.Stdout, nil))
	logger := slog.New(contextHandler)
	loggingMiddleware := middleware.NewLogging(logger)

	// Initialize any required models
	stockSource := models.NewStockTickerSource(tickerConfiguration.AlphaVantageAPIKey)

	// Initialize all handlers
	stockTickerHandler := handlers.NewStockTickerHandler(stockSource, tickerConfiguration.DefaultNumberOfDays, "MSFT", logger)
	livenessHandler := handlers.NewLivenessHandler()
	readinessHandler := handlers.NewReadinessHandler()

	// Attach logging middleware
	loggedTickerHandler := loggingMiddleware.AttachMiddleware(stockTickerHandler.GetStockData)
	loggedLivenessHandler := loggingMiddleware.AttachMiddleware(livenessHandler.Liveness)
	loggedReadinessHandler := loggingMiddleware.AttachMiddleware(readinessHandler.Readiness)

	// Attach handlers
	router := http.NewServeMux()

	router.Handle("GET /ticker", loggedTickerHandler)
	router.Handle("GET /livez", loggedLivenessHandler)
	router.Handle("GET /readyz", loggedReadinessHandler)

	logger.Info("Starting stock-ticker")

	http.ListenAndServe("0.0.0.0:8000", router)

}
