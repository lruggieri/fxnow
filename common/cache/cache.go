package cache

import (
	"context"
	"time"
)

// Cache is an interface of a service that can cache the data
//
// It is used to cache the result of a query or any data that need to reduce
// computation time
type Cache interface {
	// Set sets the value to the cache
	Set(
		ctx context.Context,
		key string,
		value interface{},
		expiration time.Duration,
	) (err error)

	// Get gets the value from the cache
	Get(
		ctx context.Context,
		key string,
		value interface{},
	) (exist bool, err error)

	// Remove removes the value from the cache
	Remove(ctx context.Context, key string) error
}
