package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
	GetAll() map[string]interface{}
	Delete(key string)
	DeleteAll()
}

type RedisCache struct{}

func NewCache(addr, password string, db int) Cache {
	InitRedis(addr, password, db)
	return &RedisCache{}
}

func InitRedis(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Printf("Error setting key: %s\n", err)
	}
}

func (c *RedisCache) Get(key string) (interface{}, bool) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		fmt.Printf("Error getting key: %s\n", err)
		return nil, false
	}
	return val, true
}

func (c *RedisCache) GetAll() map[string]interface{} {
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Printf("Error getting all keys: %s\n", err)
		return nil
	}
	result := make(map[string]interface{})
	for _, key := range keys {
		val, _ := rdb.Get(ctx, key).Result()
		result[key] = val
	}
	return result
}

func (c *RedisCache) Delete(key string) {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("Error deleting key: %s\n", err)
	}
}

func (c *RedisCache) DeleteAll() {
	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		fmt.Printf("Error deleting all keys: %s\n", err)
	}
}
