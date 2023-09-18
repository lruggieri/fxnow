package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/lruggieri/fxnow/common/cache"
	"github.com/lruggieri/fxnow/common/clock"
	cError "github.com/lruggieri/fxnow/common/error"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store"
	"github.com/lruggieri/fxnow/common/util"
)

const (
	RateLimitDuration  = time.Minute
	RateLimitMaxUsages = 2
)

type Logic interface {
	GetRate(context.Context, GetRateRequest) (*GetRateResponse, error)
}

type Impl struct {
	Store store.Store
	Cache cache.Cache
	Clock clock.Clock
}

func (i *Impl) GetRate(ctx context.Context, req GetRateRequest) (*GetRateResponse, error) {
	apiKeyID := GetAPIKeyIDFromContext(ctx)
	if len(apiKeyID) == 0 {
		return nil, errors.Wrap(cError.ErrNotAuthorized, "API key not set")
	}

	// check if API Key is valid by taking it from cache
	var cak cache.CachedAPIKey

	exist, err := i.Cache.Get(ctx, cache.GenerateCacheKeyAPIKey(apiKeyID), &cak)
	if err != nil {
		return nil, err
	}

	// if it's not found, search in DB
	if !exist {
		var res *store.GetAPIKeyResponse

		// fetch it from DB
		res, err = i.Store.GetAPIKey(ctx, store.GetAPIKeyRequest{APIKeyID: apiKeyID})
		if err != nil {
			return nil, err
		}

		cak = cache.CachedAPIKey{
			APIKeyID: apiKeyID,
			Type:     res.APIKey.Type.Uint8(),
		}
	}

	timeFrameOfInterest := i.Clock.Now().Add(-RateLimitDuration)

	// based on the API key Type, perform rate limiting
	if cak.Type == model.APIKeyTypeLimited.Uint8() &&
		!APIKeyUsagesWithinAllowedRange(cak.Usages, timeFrameOfInterest, RateLimitMaxUsages) {
		return nil, cError.ErrTooManyRequests
	}

	responseRates, err := i.fetchRates(ctx, req.Pairs)
	if err != nil {
		return nil, err
	}

	// remove useless usages
	cak.Usages = RemoveUsages(cak.Usages, timeFrameOfInterest.Unix())
	// add this usage
	cak.Usages = append(cak.Usages, cache.CachedAPIKeyUsage{Timestamp: i.Clock.Now().Unix()})
	// update cached value
	if err = i.Cache.Set(ctx, cache.GenerateCacheKeyAPIKey(apiKeyID), cak, cache.MaxCacheLifetime); err != nil {
		return nil, err
	}

	return &GetRateResponse{
		Rates: responseRates,
	}, nil
}

func (i *Impl) fetchRates(ctx context.Context, pairs []string) ([]GetRateResponseRate, error) {
	responseRates := make([]GetRateResponseRate, 0, len(pairs))

	for _, pair := range pairs {
		from, to := util.CurrenciesFromPair(pair)

		var cachedRate cache.CachedRate

		exist, err := i.Cache.Get(ctx, cache.GenerateCacheKeyRate(from, to), &cachedRate)
		if err != nil {
			return nil, err
		}

		if !exist {
			logger.WithField("pair", pair).Error("rate for pair not found")
			return nil, errors.Wrap(cError.ErrNotFound, fmt.Sprintf("rate for pair '%s' not found", pair))
		}

		responseRates = append(responseRates, GetRateResponseRate{
			Pair:      pair,
			Rate:      cachedRate.Rate,
			Timestamp: cachedRate.Timestamp,
		})
	}

	return responseRates, nil
}

func APIKeyUsagesWithinAllowedRange(usages []cache.CachedAPIKeyUsage, fromTime time.Time, maxAllowedUsages int) bool {
	usagesWithinTimeRange := 0

	for _, usage := range usages {
		if !time.Unix(usage.Timestamp, 0).Before(fromTime) {
			usagesWithinTimeRange++
		}
	}

	return usagesWithinTimeRange < maxAllowedUsages
}

// RemoveUsages : in order to spare space in the cache and make searches faster, let's remove Usages that we already
// know won't provide any useful information
func RemoveUsages(usages []cache.CachedAPIKeyUsage, notBefore int64) []cache.CachedAPIKeyUsage {
	newUsages := make([]cache.CachedAPIKeyUsage, 0, len(usages))

	for _, usage := range usages {
		if usage.Timestamp >= notBefore {
			newUsages = append(newUsages, usage)
		}
	}

	return newUsages
}
