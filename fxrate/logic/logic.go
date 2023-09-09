package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/lruggieri/fxnow/common/cache"
	"github.com/lruggieri/fxnow/common/clock"
	cError "github.com/lruggieri/fxnow/common/error"
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store"
)

const (
	RateLimitDuration  = time.Minute
	RateLimitMaxUsages = 2
	MaxCacheLifetime   = 10 * time.Minute

	// Cache prefixes
	CachePrefixAPIKey = "api_key"
	CachePrefixRate   = "rate"
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
		return nil, errors.Wrap(cError.NotAuthorized, "API key not set")
	}

	// check if API Key is valid by taking it from cache
	var cak CachedAPIKey
	exist, err := i.Cache.Get(ctx, GenerateCacheKeyAPIKey(apiKeyID), &cak)
	if err != nil {
		return nil, err
	}

	// if it's not found, search in DB
	if !exist {
		// fetch it from DB
		res, err := i.Store.GetAPIKey(ctx, store.GetAPIKeyRequest{APIKeyID: apiKeyID})
		if err != nil {
			return nil, err
		}

		cak = CachedAPIKey{
			APIKeyID: apiKeyID,
			Type:     res.APIKey.Type.Uint8(),
		}
	}

	timeFrameOfInterest := time.Now().Add(-RateLimitDuration)

	// based on the API key Type, perform rate limiting
	if cak.Type == model.APIKeyTypeLimited.Uint8() &&
		!APIKeyUsagesWithinAllowedRange(cak.Usages, timeFrameOfInterest, RateLimitMaxUsages) {
		return nil, cError.TooManyRequests
	}

	// fetch rate
	var cachedRate CachedRate

	exist, err = i.Cache.Get(ctx, GenerateCacheKeyRate(req.FromCurrency, req.ToCurrency), &cachedRate)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.Wrap(cError.NotFound, "rate for this pair not found")
	}

	// remove useless usages
	cak.Usages = RemoveUsages(cak.Usages, timeFrameOfInterest.Unix())
	// add this usage
	cak.Usages = append(cak.Usages, CachedAPIKeyUsage{Timestamp: i.Clock.Now().Unix()})
	// update cached value
	if err = i.Cache.Set(ctx, GenerateCacheKeyAPIKey(apiKeyID), cak, MaxCacheLifetime); err != nil {
		return nil, err
	}

	return &GetRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         cachedRate.Rate,
		Timestamp:    cachedRate.Timestamp,
	}, nil
}

func APIKeyUsagesWithinAllowedRange(usages []CachedAPIKeyUsage, fromTime time.Time, maxAllowedUsages int) bool {
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
func RemoveUsages(usages []CachedAPIKeyUsage, notBefore int64) []CachedAPIKeyUsage {
	newUsages := make([]CachedAPIKeyUsage, 0, len(usages))

	for _, usage := range usages {
		if usage.Timestamp >= notBefore {
			newUsages = append(newUsages, usage)
		}
	}

	return newUsages
}

func GenerateCacheKeyAPIKey(apiKeyID string) string {
	return fmt.Sprintf("%s_%s", CachePrefixAPIKey, apiKeyID)
}

func GenerateCacheKeyRate(fromCurrency, toCurrency string) string {
	return fmt.Sprintf("%s_%s_%s",
		CachePrefixRate,
		strings.ToLower(fromCurrency),
		strings.ToLower(toCurrency),
	)
}
