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
	"github.com/lruggieri/fxnow/common/logger"
)

const (
	APIURL           = "https://api.fastforex.io"
	FetchOneEndpoint = "fetch-one"
)

type Client interface {
	FetchOne(context.Context, FetchOneRequest) (*FetchOneResponse, error)
}

type FetchOneRequest struct {
	From, To string
}

type FetchOneResponse struct {
	From      string
	To        string
	Rate      float64
	Timestamp int64 // unix (s)
}

type clientFetchOneResponse struct {
	Error   string             `json:"error"`
	Base    string             `json:"base"`
	Result  map[string]float64 `json:"result"`
	Updated string             `json:"updated"`
	Ms      int                `json:"ms"`
}

type Impl struct {
	APIKey     string
	HTTPClient httpclient.Client
}

func (i *Impl) FetchOne(ctx context.Context, req FetchOneRequest) (*FetchOneResponse, error) {
	url := fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
		APIURL,
		FetchOneEndpoint,
		req.From,
		req.To,
		i.APIKey,
	)

	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	httpReq.Header.Add("accept", "application/json")

	httpResp, err := i.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}

	var response clientFetchOneResponse
	if err = json.Unmarshal(body, &response); err != nil {
		logger.WithError(err).WithField("response", string(body)).Error("invalid service response")

		return nil, errors.Wrap(err, "invalid service response")
	}

	rate, ok := response.Result[strings.ToUpper(req.To)]
	if !ok {
		logger.WithError(err).WithField("response", string(body)).Error("result is incorrect")

		return nil, errors.New("result is incorrect")
	}

	timestamp, err := time.Parse(time.DateTime, response.Updated)
	if err != nil {
		logger.WithError(err).WithField("response", string(body)).Error("invalid timestamp")

		return nil, errors.Wrap(err, "invalid timestamp")
	}

	return &FetchOneResponse{
		From:      req.From,
		To:        req.To,
		Rate:      rate,
		Timestamp: timestamp.UTC().Unix(),
	}, nil
}
