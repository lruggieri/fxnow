package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config redis config
//
// refer to this: https://github.com/redis/go-redis/blob/8db53fadf688518c2e0190619c0c4d70f0bab6a9/sentinel.go
type Config struct {
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Addrs []string `validate:"required"`

	ClusterMode bool

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string

	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int

	// Common options.

	Dialer    func(ctx context.Context, network, addr string) (net.Conn, error) `json:"-"`
	OnConnect func(ctx context.Context, cn *redis.Conn) error                   `json:"-"`

	Username         string
	Password         string
	SentinelUsername string
	SentinelPassword string

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout           time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	ContextTimeoutEnabled bool

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool

	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int // MIN_IDLE_CONNS
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration

	TLSConfig *tls.Config `json:"-"`

	// Only cluster clients.

	MaxRedirects   int
	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool

	// The sentinel master name.
	// Only failover clients.
	MasterName string
}

func (c *Config) ToRedisOptions() *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:                 c.Addrs,
		ClientName:            c.ClientName,
		DB:                    c.DB,
		Dialer:                c.Dialer,
		OnConnect:             c.OnConnect,
		Username:              c.Username,
		Password:              c.Password,
		SentinelUsername:      c.SentinelUsername,
		SentinelPassword:      c.SentinelPassword,
		MaxRetries:            c.MaxRetries,
		MinRetryBackoff:       c.MinRetryBackoff,
		MaxRetryBackoff:       c.MaxRetryBackoff,
		DialTimeout:           c.DialTimeout,
		ReadTimeout:           c.ReadTimeout,
		WriteTimeout:          c.WriteTimeout,
		ContextTimeoutEnabled: c.ContextTimeoutEnabled,
		PoolFIFO:              c.PoolFIFO,
		PoolSize:              c.PoolSize,
		PoolTimeout:           c.PoolTimeout,
		MinIdleConns:          c.MinIdleConns,
		MaxIdleConns:          c.MaxIdleConns,
		ConnMaxIdleTime:       c.ConnMaxIdleTime,
		ConnMaxLifetime:       c.ConnMaxLifetime,
		TLSConfig:             c.TLSConfig,
		MaxRedirects:          c.MaxRedirects,
		ReadOnly:              c.ReadOnly,
		RouteByLatency:        c.RouteByLatency,
		RouteRandomly:         c.RouteRandomly,
		MasterName:            c.MasterName,
	}
}
