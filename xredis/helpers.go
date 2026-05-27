package xredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Ping verifies Redis connectivity.
func Ping(ctx context.Context, client redis.UniversalClient) error {
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("xredis: ping failed: %w", err)
	}
	return nil
}

// Close shuts down the Redis client, releasing all connections.
func Close(client redis.UniversalClient) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("xredis: close failed: %w", err)
	}
	return nil
}
