// Package redis opens the shared go-redis client used for session cache, rate limit
// buckets, pubsub fanout, presence, and idempotency lookup.
package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Options builds client options from host:port + password + db, or from a
// redis:// / rediss:// URL in addr (hosts such as Render hand out a
// connection string; DD-239), which then carries the password and db.
func Options(addr, password string, dbIdx int) (*redis.Options, error) {
	o := &redis.Options{Addr: addr, Password: password, DB: dbIdx}
	if strings.HasPrefix(addr, "redis://") || strings.HasPrefix(addr, "rediss://") {
		parsed, err := redis.ParseURL(addr)
		if err != nil {
			return nil, fmt.Errorf("redis url: %w", err)
		}
		o = parsed
	}
	o.PoolSize = 50
	o.MinIdleConns = 5
	o.DialTimeout = 5 * time.Second
	o.ReadTimeout = 3 * time.Second
	o.WriteTimeout = 3 * time.Second
	return o, nil
}

// Connect creates the client, applies pool sizing, and pings before returning.
func Connect(ctx context.Context, addr, password string, dbIdx int) (*redis.Client, error) {
	opts, err := Options(addr, password, dbIdx)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return client, nil
}
