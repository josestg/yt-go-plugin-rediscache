package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/josestg/yt-go-plugin/cache"
	"github.com/redis/go-redis/v9"
)

// RedisCache is a cache implementation that uses Redis.
type RedisCache struct {
	log    *slog.Logger
	client *redis.Client
}

// Factory is the symbol the plugin loader will try to load. It must implement the cache.Factory signature.
var Factory cache.Factory = New

// New creates a new RedisCache instance.
func New(log *slog.Logger) (cache.Cache, error) {
	log.Info("[plugin/rediscache] loaded")
	db, err := strconv.Atoi(cmp.Or(os.Getenv("REDIS_DB"), "0"))
	if err != nil {
		return nil, fmt.Errorf("parse redis db: %w", err)
	}

	c := &RedisCache{
		log: log,
		client: redis.NewClient(&redis.Options{
			Addr:     cmp.Or(os.Getenv("REDIS_ADDR"), "localhost:6379"),
			Password: cmp.Or(os.Getenv("REDIS_PASSWORD"), ""),
			DB:       db,
		}),
	}

	return c, nil
}

func (r *RedisCache) Set(ctx context.Context, key, val string, exp time.Duration) error {
	r.log.InfoContext(ctx, "[plugin/rediscache] set", "key", key, "val", val, "exp", exp)
	return r.client.Set(ctx, key, val, exp).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	r.log.InfoContext(ctx, "[plugin/rediscache] get", "key", key)
	res, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		r.log.InfoContext(ctx, "[plugin/rediscache] key not found", "key", key)
		return "", cache.ErrNotFound
	}
	r.log.InfoContext(ctx, "[plugin/rediscache] key found", "key", key, "val", res)
	return res, err
}
