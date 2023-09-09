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

	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"
	mockhttpclient "github.com/lruggieri/fxnow/common/mock/client/httpclient"
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
		req FetchOneRequest
	}

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			res *FetchOneResponse,
			err error,
		)
	}{
		{
			name: "error-wrong-response",
			args: args{
				ctx: context.Background(),
				req: FetchOneRequest{
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
			assertion: func(t *testing.T, res *FetchOneResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "error-invalid-to-currency",
			args: args{
				ctx: context.Background(),
				req: FetchOneRequest{
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
			assertion: func(t *testing.T, res *FetchOneResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "error-invalid-timestamp",
			args: args{
				ctx: context.Background(),
				req: FetchOneRequest{
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
			assertion: func(t *testing.T, res *FetchOneResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "happy-path",
			args: args{
				ctx: context.Background(),
				req: FetchOneRequest{
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
			assertion: func(t *testing.T, res *FetchOneResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &FetchOneResponse{
					From:      "USD",
					To:        "JPY",
					Rate:      42.42,
					Timestamp: now.Unix(),
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

			l := Impl{
				APIKey:     apiKey,
				HTTPClient: d.httpClient,
			}

			tc.mock(tc.args, d)

			res, err := l.FetchOne(tc.args.ctx, tt.args.req)

			tc.assertion(t, res, err)
		})
	}
}
