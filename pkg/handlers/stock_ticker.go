package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/naspinall/stock-ticker/pkg/models"
)

const timeParsingLayout = "2006-01-02"
const daysQueryKey = "days"
const stockSourceQueryKey = "source"

// Output structs
type StockTickerEntry struct {
	Open   float64 `json:"open_price"`
	High   float64 `json:"high_price"`
	Low    float64 `json:"low_price"`
	Close  float64 `json:"close_price"`
	Volume int     `json:"total_volume"`
	Date   string  `json:"date"`
}

type StockTickerResponse struct {
	AverageClosingPrice float64            `json:"average_closing_price"`
	TimeSeriesData      []StockTickerEntry `json:"time_series_data"`
	Timezone            string             `json:"timezone"`
	Symbol              string             `json:"symbol"`
	LastRefreshed       string             `json:"last_refreshed"`
}

type StockTickerHandler struct {
	defaultNumberOfDays int
	defaultStockSource  string
	stockTickerSource   *models.StockTickerSource
	logger              *slog.Logger
}

func NewStockTickerHandler(stockTickerSource *models.StockTickerSource, defaultNumberOfDays int, defaultStockSource string, logger *slog.Logger) *StockTickerHandler {
	return &StockTickerHandler{
		defaultNumberOfDays: defaultNumberOfDays,
		defaultStockSource:  defaultStockSource,
		stockTickerSource:   stockTickerSource,
		logger:              logger,
	}
}

// Convenience function to parse an int query string
func (sth *StockTickerHandler) parseIntQuery(values url.Values, key string) (int, error) {

	valueString := values.Get(key)
	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		return -1, nil
	}

	return valueInt, nil

}

// Convenience function to parse out query string for request
func (sth *StockTickerHandler) parseQuery(r *http.Request) (int, string, error) {

	numberOfDays := sth.defaultNumberOfDays
	stockSource := sth.defaultStockSource

	queryValues := r.URL.Query()
	if queryValues.Has(daysQueryKey) {
		parsedNumberOfDays, err := sth.parseIntQuery(queryValues, daysQueryKey)
		if err != nil {
			sth.logger.ErrorContext(r.Context(), "invalid `days` query string", "error", err)
			return 0, "", err
		}

		numberOfDays = parsedNumberOfDays
	}

	if queryValues.Has(stockSourceQueryKey) {
		stockSource = queryValues.Get(stockSourceQueryKey)
	}

	return numberOfDays, stockSource, nil
}

// Calculate close price average
func calculateAverageClosingPrice(stockTickerEntries []StockTickerEntry, dayCount int) float64 {

	var closePriceRunningTotal float64 = 0

	for _, stockTickerEntry := range stockTickerEntries {
		closePriceRunningTotal += stockTickerEntry.Close
	}

	return closePriceRunningTotal / float64(dayCount)
}

func (sth *StockTickerHandler) GetStockData(w http.ResponseWriter, r *http.Request) {

	// Parse out optional query parameters
	dayCount, stockSource, err := sth.parseQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get stock information
	stockSourceResponse, err := sth.stockTickerSource.GetStockSource(stockSource)
	if err != nil {
		sth.logger.ErrorContext(r.Context(), "cannot get stock source information", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse out response stings into date objects for sorting
	timeSeriesDateStrings := []time.Time{}
	for date := range stockSourceResponse.TimeSeriesDaily {

		dateString, err := time.Parse(timeParsingLayout, date)
		if err != nil {
			sth.logger.ErrorContext(r.Context(), "cannot parse source date time string", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		timeSeriesDateStrings = append(timeSeriesDateStrings, dateString)

	}

	// Sort by date descending
	sort.SliceStable(timeSeriesDateStrings, func(left, right int) bool {
		return timeSeriesDateStrings[left].After(timeSeriesDateStrings[right])
	})

	stockTickerEntries := []StockTickerEntry{}

	for index := 0; index < dayCount; index++ {

		dateString := timeSeriesDateStrings[index].Format(timeParsingLayout)

		dayData := stockSourceResponse.TimeSeriesDaily[dateString]

		stockTickerEntry := StockTickerEntry{
			Open:   dayData.Open,
			High:   dayData.High,
			Low:    dayData.Low,
			Close:  dayData.Close,
			Volume: dayData.Volume,
			Date:   dateString,
		}

		stockTickerEntries = append(stockTickerEntries, stockTickerEntry)
	}

	averageClosingPrice := calculateAverageClosingPrice(stockTickerEntries, dayCount)

	stockTickerResponse := StockTickerResponse{
		AverageClosingPrice: averageClosingPrice,
		TimeSeriesData:      stockTickerEntries,
		Timezone:            stockSourceResponse.Metadata.TimeZone,
		Symbol:              stockSourceResponse.Metadata.Symbol,
		LastRefreshed:       stockSourceResponse.Metadata.LastRefreshed,
	}

	err = json.NewEncoder(w).Encode(&stockTickerResponse)
	if err != nil {
		sth.logger.ErrorContext(r.Context(), "cannot encode output", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
