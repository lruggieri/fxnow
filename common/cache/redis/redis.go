package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Cacher struct {
	Client UniversalClient
}

// Get implements cache.Cacher
func (c *Cacher) Get(ctx context.Context, key string, value interface{}) (exist bool, err error) {
	data, err := c.Client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if len(data) == 0 || err == redis.Nil {
		return false, nil
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return false, errors.Wrap(err, "unmarshal data")
	}

	return true, nil
}

// Remove implements cache.Cacher
func (c *Cacher) Remove(ctx context.Context, key string) error {
	err := c.Client.Del(ctx, key).Err()

	return errors.Wrap(err, "cannot remove key from redis")
}

// Set implements cache.Cacher
func (c *Cacher) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Client.Set(ctx, key, string(data), expiration).Err()

	return errors.Wrap(err, "cannot set key to redis")
}

type UniversalClient interface {
	redis.UniversalClient
}

type Configure func(opt *Config)

func NewClient(c Config, configures ...Configure) UniversalClient {
	for _, configure := range configures {
		configure(&c)
	}

	uniOpts := *c.ToRedisOptions()

	if c.ClusterMode {
		return redis.NewClusterClient(uniOpts.Cluster())
	}

	return redis.NewUniversalClient(&uniOpts)
}
