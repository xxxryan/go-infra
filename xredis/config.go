package xredis

import (
	"crypto/tls"
	"time"
)

type Config struct {
	// Standalone mode: single Redis address (default "localhost:6379").
	Addr string

	// Sentinel mode: set MasterName and SentinelAddrs to enable.
	MasterName    string
	SentinelAddrs []string

	// Cluster mode: set ClusterAddrs to enable.
	ClusterAddrs []string

	// Common options.
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
	TLSConfig    *tls.Config
}

func (c Config) mode() string {
	switch {
	case len(c.ClusterAddrs) > 0:
		return "cluster"
	case c.MasterName != "":
		return "sentinel"
	default:
		return "standalone"
	}
}
