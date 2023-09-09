package logic

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cError "github.com/lruggieri/fxnow/common/error"
	mockcache "github.com/lruggieri/fxnow/common/mock/cache"
	mockclock "github.com/lruggieri/fxnow/common/mock/clock"
	mockstore "github.com/lruggieri/fxnow/common/mock/store"
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store"
)

func TestLogicGetRate(t *testing.T) {
	testErr := errors.New("error")
	now := time.Now()

	type deps struct {
		store *mockstore.Store
		cache *mockcache.Cache
		clock *mockclock.Clock
	}

	type args struct {
		ctx context.Context
		req GetRateRequest
	}

	apiKey := "api_key"
	apiKeyCtx := context.WithValue(context.Background(), ContextKeyAPIKey, apiKey)

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			res *GetRateResponse,
			err error,
		)
	}{
		{
			name: "error-no-api-key-set",
			args: args{
				ctx: context.Background(),
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, cError.ErrNotAuthorized)
			},
		},
		{
			name: "error-get-cache-api-key",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).Return(false, testErr).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-get-cache-rate",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
						},
					}))
					return true, nil
				}).Once()

				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyRate(args.req.FromCurrency, args.req.ToCurrency),
					mock.AnythingOfType("*logic.CachedRate"),
				).Return(false, testErr).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-db-fetch",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).Return(false, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{APIKeyID: apiKey}).
					Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-set-cache",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
						},
					}))
					return true, nil
				}).Once()

				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyRate(args.req.FromCurrency, args.req.ToCurrency),
					mock.AnythingOfType("*logic.CachedRate"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedRate{
						Rate:      42.42,
						Timestamp: now.Unix(),
					}))
					return true, nil
				}).Once()

				d.clock.EXPECT().Now().Return(now).Once()

				d.cache.EXPECT().Set(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
							{Timestamp: now.Unix()},
						},
					},
					MaxCacheLifetime,
				).Return(testErr).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-too-many-requests",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
							{Timestamp: now.Unix()},
						},
					}))
					return true, nil
				}).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, cError.ErrTooManyRequests)
			},
		},
		{
			name: "happy-path-limited-within-range",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
						},
					}))
					return true, nil
				}).Once()

				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyRate(args.req.FromCurrency, args.req.ToCurrency),
					mock.AnythingOfType("*logic.CachedRate"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedRate{
						Rate:      42.42,
						Timestamp: now.Unix(),
					}))
					return true, nil
				}).Once()

				d.clock.EXPECT().Now().Return(now).Once()

				d.cache.EXPECT().Set(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
							{Timestamp: now.Unix()},
						},
					},
					MaxCacheLifetime,
				).Return(nil).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &GetRateResponse{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
					Rate:         42.42,
					Timestamp:    now.Unix(),
				}, res)
			},
		},
		{
			name: "happy-path-no-api-key-cache",
			args: args{
				ctx: apiKeyCtx,
				req: GetRateRequest{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
				},
			},
			mock: func(args args, d deps) {
				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					mock.AnythingOfType("*logic.CachedAPIKey"),
				).Return(false, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{APIKeyID: apiKey}).
					Return(&store.GetAPIKeyResponse{APIKey: &model.APIKey{
						APIKeyID: "api_key",
						Type:     model.APIKeyTypeLimited,
					}}, nil).Once()

				d.cache.EXPECT().Get(
					args.ctx,
					GenerateCacheKeyRate(args.req.FromCurrency, args.req.ToCurrency),
					mock.AnythingOfType("*logic.CachedRate"),
				).RunAndReturn(func(ctx context.Context, s string, i interface{}) (bool, error) {
					reflect.ValueOf(i).Elem().Set(reflect.ValueOf(CachedRate{
						Rate:      42.42,
						Timestamp: now.Unix(),
					}))
					return true, nil
				}).Once()

				d.clock.EXPECT().Now().Return(now).Once()

				d.cache.EXPECT().Set(
					args.ctx,
					GenerateCacheKeyAPIKey(apiKey),
					CachedAPIKey{
						APIKeyID: apiKey,
						Type:     model.APIKeyTypeLimited.Uint8(),
						Usages: []CachedAPIKeyUsage{
							{Timestamp: now.Unix()},
						},
					},
					MaxCacheLifetime,
				).Return(nil).Once()
			},
			assertion: func(t *testing.T, res *GetRateResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &GetRateResponse{
					FromCurrency: "USD",
					ToCurrency:   "JPY",
					Rate:         42.42,
					Timestamp:    now.Unix(),
				}, res)
			},
		},
	}

	for _, tt := range tests {
		tc := tt // avoid loop closure issue
		t.Run(tc.name, func(t *testing.T) {
			d := deps{
				store: mockstore.NewStore(t),
				cache: mockcache.NewCache(t),
				clock: mockclock.NewClock(t),
			}

			l := Impl{
				Store: d.store,
				Cache: d.cache,
				Clock: d.clock,
			}

			tc.mock(tc.args, d)

			res, err := l.GetRate(tc.args.ctx, tt.args.req)

			tc.assertion(t, res, err)
		})
	}
}
