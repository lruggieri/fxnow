package fastforex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/lruggieri/fxnow/common/client/httpclient"
	"github.com/lruggieri/fxnow/common/fxsource"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/util"
)

const (
	APIURL           = "https://api.fastforex.io"
	FetchOneEndpoint = "fetch-one"
	FetchAllEndpoint = "fetch-all"
)

type clientFetchOneResponse struct {
	Error   string             `json:"error"`
	Base    string             `json:"base"`
	Result  map[string]float64 `json:"result"`
	Updated string             `json:"updated"`
	Ms      int                `json:"ms"`
}

type clientFetchAllResponse struct {
	Error   string             `json:"error"`
	Base    string             `json:"base"`
	Results map[string]float64 `json:"results"`
	Updated string             `json:"updated"`
	Ms      int                `json:"ms"`
}

type Client struct {
	APIKey     string
	HTTPClient httpclient.Client
}

func (i *Client) FetchRate(ctx context.Context, req fxsource.FetchRateRequest) (*fxsource.FetchRateResponse, error) {
	url := fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
		APIURL,
		FetchOneEndpoint,
		req.From,
		req.To,
		i.APIKey,
	)

	var response clientFetchOneResponse
	if err := i.fetch(ctx, url, &response); err != nil {
		return nil, errors.Wrap(err, "cannot fetch rate")
	}

	rate, ok := response.Result[strings.ToUpper(req.To)]
	if !ok {
		logger.WithField("response", response).Error("result is incorrect")

		return nil, errors.New("result is incorrect")
	}

	timestamp, err := time.Parse(time.DateTime, response.Updated)
	if err != nil {
		logger.WithError(err).WithField("response", response).Error("invalid timestamp")

		return nil, errors.Wrap(err, "invalid timestamp")
	}

	return &fxsource.FetchRateResponse{
		Rate: fxsource.Rate{
			From:      req.From,
			To:        req.To,
			Rate:      rate,
			Timestamp: timestamp.UTC().Unix(),
		},
	}, nil
}

func (i *Client) FetchAllRates(
	ctx context.Context, req fxsource.FetchAllRatesRequest,
) (*fxsource.FetchAllRatesResponse, error) {
	currenciesToFetch := fiatCurrencies

	if len(req.Limit) > 0 {
		limitMap := util.SliceToMap(req.Limit)

		for currency := range currenciesToFetch {
			if _, ok := limitMap[currency]; !ok {
				delete(currenciesToFetch, currency)
			}
		}
	}

	rates := make([]fxsource.Rate, 0)

	for currency := range currenciesToFetch {
		if err := i.fetchAllRateCurrency(ctx, currency, func(rate fxsource.Rate) {
			rates = append(rates, rate)
		}); err != nil {
			return nil, err
		}
	}

	return &fxsource.FetchAllRatesResponse{
		Rates: rates,
	}, nil
}

func (i *Client) fetchAllRateCurrency(ctx context.Context, currency string, addRate func(fxsource.Rate)) error {
	url := fmt.Sprintf("%s/%s?from=%s&api_key=%s",
		APIURL,
		FetchAllEndpoint,
		currency,
		i.APIKey,
	)

	var response clientFetchAllResponse
	if err := i.fetch(ctx, url, &response); err != nil {
		return errors.Wrap(err, "cannot fetch rate")
	}

	timestamp, err := time.Parse(time.DateTime, response.Updated)
	if err != nil {
		logger.WithError(err).WithField("response", response).Error("invalid timestamp")

		return errors.Wrap(err, "invalid timestamp")
	}

	for toCurrency, rate := range response.Results {
		addRate(fxsource.Rate{
			From:      currency,
			To:        toCurrency,
			Rate:      rate,
			Timestamp: timestamp.UTC().Unix(),
		})
	}

	return nil
}

func (i *Client) fetch(ctx context.Context, url string, resp interface{}) error {
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	httpReq.Header.Add("accept", "application/json")

	httpResp, err := i.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		logger.WithField("status-code", httpResp.StatusCode).Error("status code != 200")

		return errors.New("status code != 200")
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}

	if err = json.Unmarshal(body, resp); err != nil {
		logger.WithError(err).WithField("response", string(body)).Error("invalid service response")

		return errors.Wrap(err, "invalid service response")
	}

	return nil
}

func NewClient(apiKey string, httpClient httpclient.Client) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: httpClient,
	}
}
