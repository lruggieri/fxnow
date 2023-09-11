package logic

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lruggieri/fxnow/common/cache"
	"github.com/lruggieri/fxnow/common/fxsource"
	mockcache "github.com/lruggieri/fxnow/common/mock/cache"
	mockfxsource "github.com/lruggieri/fxnow/common/mock/fxsource"
)

func TestImpl_fxUpdate(t *testing.T) {
	testErr := errors.New("error")
	now := time.Now()

	type deps struct {
		cache    *mockcache.Cache
		fxSource *mockfxsource.FXSource
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			err error,
		)
	}{
		{
			name: "error-fetch-rates",
			args: args{
				ctx: context.Background(),
			},
			mock: func(args args, d deps) {
				d.fxSource.EXPECT().FetchAllRates(
					args.ctx,
					fxsource.FetchAllRatesRequest{
						Limit: []string{"USD", "GBP", "EUR", "JPY", "CHF", "CAD", "AUD"},
					},
				).Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-fetch-rates",
			args: args{
				ctx: context.Background(),
			},
			mock: func(args args, d deps) {
				d.fxSource.EXPECT().FetchAllRates(
					args.ctx,
					fxsource.FetchAllRatesRequest{
						Limit: []string{"USD", "GBP", "EUR", "JPY", "CHF", "CAD", "AUD"},
					},
				).Return(&fxsource.FetchAllRatesResponse{
					Rates: []fxsource.Rate{
						{
							From:      "USD",
							To:        "JPY",
							Rate:      42.42,
							Timestamp: now.Unix(),
						},
						{
							From:      "EUR",
							To:        "CAD",
							Rate:      42.43,
							Timestamp: now.Unix(),
						},
					},
				}, nil).Once()

				d.cache.EXPECT().Set(
					args.ctx,
					cache.GenerateCacheKeyRate("USD", "JPY"),
					cache.CachedRate{
						Rate:      42.42,
						Timestamp: now.Unix(),
					},
					cache.MaxCacheLifetime,
				).Return(testErr).Once()
			},
			assertion: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "happy-path",
			args: args{
				ctx: context.Background(),
			},
			mock: func(args args, d deps) {
				d.fxSource.EXPECT().FetchAllRates(
					args.ctx,
					fxsource.FetchAllRatesRequest{
						Limit: []string{"USD", "GBP", "EUR", "JPY", "CHF", "CAD", "AUD"},
					},
				).Return(&fxsource.FetchAllRatesResponse{
					Rates: []fxsource.Rate{
						{
							From:      "USD",
							To:        "JPY",
							Rate:      42.42,
							Timestamp: now.Unix(),
						},
						{
							From:      "EUR",
							To:        "CAD",
							Rate:      42.43,
							Timestamp: now.Unix(),
						},
					},
				}, nil).Once()

				d.cache.EXPECT().Set(
					args.ctx,
					cache.GenerateCacheKeyRate("USD", "JPY"),
					cache.CachedRate{
						Rate:      42.42,
						Timestamp: now.Unix(),
					},
					cache.MaxCacheLifetime,
				).Return(nil).Once()

				d.cache.EXPECT().Set(
					args.ctx,
					cache.GenerateCacheKeyRate("EUR", "CAD"),
					cache.CachedRate{
						Rate:      42.43,
						Timestamp: now.Unix(),
					},
					cache.MaxCacheLifetime,
				).Return(nil).Once()
			},
			assertion: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		tc := tt // avoid loop closure issue
		t.Run(tc.name, func(t *testing.T) {
			d := deps{
				cache:    mockcache.NewCache(t),
				fxSource: mockfxsource.NewFXSource(t),
			}

			l := Impl{
				Cache:    d.cache,
				FXSource: d.fxSource,
			}

			tc.mock(tc.args, d)

			err := l.fxUpdate(tc.args.ctx)

			tc.assertion(t, err)
		})
	}
}
