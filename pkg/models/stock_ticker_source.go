package models

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const stockTickerAPIHost = "www.alphavantage.co"

type StockTickerSource struct {
	apiKey string
}

func NewStockTickerSource(apiKey string) *StockTickerSource {
	return &StockTickerSource{
		apiKey: apiKey,
	}
}

type StockMetadata struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type TimeSeriesEntry struct {
	Open   float64 `json:"1. open,string"`
	High   float64 `json:"2. high,string"`
	Low    float64 `json:"3. low,string"`
	Close  float64 `json:"4. close,string"`
	Volume int     `json:"5. volume,string"`
}

type StockSourceResponse struct {
	Metadata        StockMetadata              `json:"Meta Data"`
	TimeSeriesDaily map[string]TimeSeriesEntry `json:"Time Series (Daily)"`
}

func (sts *StockTickerSource) GetStockSource(symbol string) (*StockSourceResponse, error) {

	queryParams := url.Values{}
	queryParams.Add("apikey", sts.apiKey)
	queryParams.Add("function", "TIME_SERIES_DAILY")
	queryParams.Add("symbol", symbol)

	tickerSourceURL := url.URL{
		Scheme:   "https",
		Host:     stockTickerAPIHost,
		Path:     "/query",
		RawQuery: queryParams.Encode(),
	}

	response, err := http.Get(tickerSourceURL.String())
	if err != nil {
		return nil, err
	}

	stockSourceResponse := &StockSourceResponse{}

	err = json.NewDecoder(response.Body).Decode(stockSourceResponse)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return stockSourceResponse, nil
}
