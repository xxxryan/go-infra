package xredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Open creates a Redis client based on the config mode and verifies connectivity.
func Open(cfg Config) (redis.UniversalClient, error) {
	client := newClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout(cfg))
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("xredis: ping failed: %w", err)
	}
	return client, nil
}

func newClient(cfg Config) redis.UniversalClient {
	switch cfg.mode() {
	case "cluster":
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.ClusterAddrs,
			Password:     cfg.Password,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			TLSConfig:    cfg.TLSConfig,
		})
	case "sentinel":
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			SentinelAddrs: cfg.SentinelAddrs,
			Password:      cfg.Password,
			DB:            cfg.DB,
			DialTimeout:   cfg.DialTimeout,
			ReadTimeout:   cfg.ReadTimeout,
			WriteTimeout:  cfg.WriteTimeout,
			PoolSize:      cfg.PoolSize,
			MinIdleConns:  cfg.MinIdleConns,
			TLSConfig:     cfg.TLSConfig,
		})
	default:
		addr := cfg.Addr
		if addr == "" {
			addr = "localhost:6379"
		}
		return redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			TLSConfig:    cfg.TLSConfig,
		})
	}
}

func dialTimeout(cfg Config) time.Duration {
	if cfg.DialTimeout > 0 {
		return cfg.DialTimeout
	}
	return 5 * time.Second
}
