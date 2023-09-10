package fxsource

import "context"

type FXSource interface {
	FetchRate(context.Context, FetchRateRequest) (*FetchRateResponse, error)
	FetchAllRates(context.Context, FetchAllRatesRequest) (*FetchAllRatesResponse, error)
}

type FetchRateRequest struct {
	From string
	To   string
}

type FetchRateResponse struct {
	Rate
}

type Rate struct {
	From      string
	To        string
	Rate      float64
	Timestamp int64 // unix (s)
}

type FetchAllRatesRequest struct {
	// Limit: Optional. List of currencies we want to fetch rates for.
	Limit []string
}

type FetchAllRatesResponse struct {
	Rates []Rate
}
