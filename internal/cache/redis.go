package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) Del(key string) error {
	return r.client.Del(ctx, key).Err()
}
