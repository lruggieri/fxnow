package logic

import (
	"context"
	"time"

	"github.com/lruggieri/fxnow/common/cache"
	"github.com/lruggieri/fxnow/common/fxsource"
	"github.com/lruggieri/fxnow/common/logger"
)

const (
	UpdateTicker = 20 * time.Second
)

type Logic interface {
	StartFXUpdate(ctx context.Context)
}

type Impl struct {
	Cache    cache.Cache
	FXSource fxsource.FXSource
}

func (i *Impl) StartFXUpdate(ctx context.Context) {
	ticker := time.NewTicker(UpdateTicker)

	logger.Info("starting FX update loop with %s tick loops", UpdateTicker.String())

	if err := i.fxUpdate(ctx); err != nil {
		// TODO better management with notifications
		logger.WithError(err).Error("fx update error")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := i.fxUpdate(ctx); err != nil {
				// TODO better management with notifications
				logger.WithError(err).Error("fx update error")
			}
		}
	}
}

func (i *Impl) fxUpdate(ctx context.Context) error {
	rates, err := i.FXSource.FetchAllRates(ctx, fxsource.FetchAllRatesRequest{
		Limit: []string{"USD", "GBP", "EUR", "JPY", "CHF", "CAD", "AUD"},
	})
	if err != nil {
		return err
	}

	for _, rate := range rates.Rates {
		if err = i.Cache.Set(ctx, cache.GenerateCacheKeyRate(rate.From, rate.To), cache.CachedRate{
			Rate:      rate.Rate,
			Timestamp: rate.Timestamp,
		}, cache.MaxCacheLifetime); err != nil {
			return err
		}
	}

	return nil
}
