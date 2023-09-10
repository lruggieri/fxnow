package fastforex

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lruggieri/fxnow/common/fxsource"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"
	mockhttpclient "github.com/lruggieri/fxnow/common/mock/client/httpclient"
	"github.com/lruggieri/fxnow/common/util"
)

func TestImpl_FetchOne(t *testing.T) {
	logger.InitLogger(zap.New(zap.Config{Development: false}))
	// testErr := errors.New("error")
	apiKey := "api-key"
	now := time.Now().UTC()
	nowFormatted := now.Format(time.DateTime)

	type deps struct {
		httpClient *mockhttpclient.Client
	}

	type args struct {
		ctx context.Context
		req fxsource.FetchRateRequest
	}

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			res *fxsource.FetchRateResponse,
			err error,
		)
	}{
		{
			name: "error-wrong-response",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchRateRequest{
					From: "USD",
					To:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
						APIURL,
						FetchOneEndpoint,
						args.req.From,
						args.req.To,
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`wrong response`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchRateResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "error-wrong-status-code",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchRateRequest{
					From: "USD",
					To:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
						APIURL,
						FetchOneEndpoint,
						args.req.From,
						args.req.To,
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader(``)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchRateResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "error-invalid-to-currency",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchRateRequest{
					From: "USD",
					To:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
						APIURL,
						FetchOneEndpoint,
						args.req.From,
						args.req.To,
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						  "base": "USD",
						  "result": {
							"GBP": 42.42
						  },
						  "updated": "` + nowFormatted + `",
						  "ms": 11
						}`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchRateResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "error-invalid-timestamp",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchRateRequest{
					From: "USD",
					To:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
						APIURL,
						FetchOneEndpoint,
						args.req.From,
						args.req.To,
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						  "base": "USD",
						  "result": {
							"JPY": 42.42
						  },
						  "updated": "invalid timestamp",
						  "ms": 11
						}`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchRateResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "happy-path",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchRateRequest{
					From: "USD",
					To:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&to=%s&api_key=%s",
						APIURL,
						FetchOneEndpoint,
						args.req.From,
						args.req.To,
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						  "base": "USD",
						  "result": {
							"JPY": 42.42
						  },
						  "updated": "` + nowFormatted + `",
						  "ms": 11
						}`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchRateResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &fxsource.FetchRateResponse{
					Rate: fxsource.Rate{
						From:      "USD",
						To:        "JPY",
						Rate:      42.42,
						Timestamp: now.Unix(),
					},
				}, res)
			},
		},
	}

	for _, tt := range tests {
		tc := tt // avoid loop closure issue
		t.Run(tc.name, func(t *testing.T) {
			d := deps{
				httpClient: mockhttpclient.NewClient(t),
			}

			l := Client{
				APIKey:     apiKey,
				HTTPClient: d.httpClient,
			}

			tc.mock(tc.args, d)

			res, err := l.FetchRate(tc.args.ctx, tt.args.req)

			tc.assertion(t, res, err)
		})
	}
}

func TestImpl_FetchAll(t *testing.T) {
	logger.InitLogger(zap.New(zap.Config{Development: false}))
	// testErr := errors.New("error")
	apiKey := "api-key"
	now := time.Now().UTC()
	nowFormatted := now.Format(time.DateTime)

	type deps struct {
		httpClient *mockhttpclient.Client
	}

	type args struct {
		ctx context.Context
		req fxsource.FetchAllRatesRequest
	}

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			res *fxsource.FetchAllRatesResponse,
			err error,
		)
	}{
		{
			name: "error-wrong-response",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchAllRatesRequest{
					Limit: []string{"USD"},
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&api_key=%s",
						APIURL,
						FetchAllEndpoint,
						"USD",
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`wrong response`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchAllRatesResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "error-timestamp-format",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchAllRatesRequest{
					Limit: []string{"USD"},
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&api_key=%s",
						APIURL,
						FetchAllEndpoint,
						"USD",
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						  "base": "USD",
						  "results": {
							"JPY": 42.42,
							"EUR": 42.43
						  },
						  "updated": "2006-01-02",
						  "ms": 11
						}`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchAllRatesResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "happy-path",
			args: args{
				ctx: context.Background(),
				req: fxsource.FetchAllRatesRequest{
					Limit: []string{"USD"},
				},
			},
			mock: func(args args, d deps) {
				d.httpClient.EXPECT().Do(mock.Anything).Run(func(req *http.Request) {
					assert.Equal(t, fmt.Sprintf("%s/%s?from=%s&api_key=%s",
						APIURL,
						FetchAllEndpoint,
						"USD",
						apiKey,
					), req.URL.String())
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						  "base": "USD",
						  "results": {
							"JPY": 42.42,
							"EUR": 42.43
						  },
						  "updated": "` + nowFormatted + `",
						  "ms": 11
						}`)),
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *fxsource.FetchAllRatesResponse, err error) {
				assert.Nil(t, err)
				assert.True(t, util.SliceUnorderedEqual(
					[]fxsource.Rate{
						{
							From:      "USD",
							To:        "JPY",
							Rate:      42.42,
							Timestamp: now.Unix(),
						},
						{
							From:      "USD",
							To:        "EUR",
							Rate:      42.43,
							Timestamp: now.Unix(),
						},
					},
					res.Rates,
				))
			},
		},
	}

	for _, tt := range tests {
		tc := tt // avoid loop closure issue
		t.Run(tc.name, func(t *testing.T) {
			d := deps{
				httpClient: mockhttpclient.NewClient(t),
			}

			l := Client{
				APIKey:     apiKey,
				HTTPClient: d.httpClient,
			}

			tc.mock(tc.args, d)

			res, err := l.FetchAllRates(tc.args.ctx, tt.args.req)

			tc.assertion(t, res, err)

			assert.True(t, d.httpClient.AssertExpectations(t))
		})
	}
}
