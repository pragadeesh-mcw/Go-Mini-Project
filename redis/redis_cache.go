package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	client  *redis.Client
	maxSize int
}

func NewCache(addr string, password string, db int, maxSize int) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Cache{
		client:  rdb,
		maxSize: maxSize,
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	c.updateAccessOrder(key)
	return c.evictIfNecessary()
}

func (c *Cache) Get(key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	c.updateAccessOrder(key)
	return val, nil
}

func (c *Cache) updateAccessOrder(key string) {
	// remove and readd to maintain LRU order
	c.client.LRem(ctx, "cache_keys", 0, key)
	c.client.LPush(ctx, "cache_keys", key)
}

func (c *Cache) GetAll() (map[string]string, error) {
	keys, err := c.client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	for _, key := range keys {
		val, err := c.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		values[key] = val
	}
	return values, nil
}

func (c *Cache) Delete(key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	c.client.LRem(ctx, "cache_keys", 0, key)
	return nil
}

func (c *Cache) DeleteAll() error {
	keys, err := c.client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		c.client.Del(ctx, key)
	}

	c.client.Del(ctx, "cache_keys")
	return nil
}

func (c *Cache) evictIfNecessary() error {
	size := c.client.LLen(ctx, "cache_keys").Val()
	if size > int64(c.maxSize) {
		excess := size - int64(c.maxSize)
		for i := int64(0); i < excess; i++ {
			key := c.client.RPop(ctx, "cache_keys").Val()
			c.client.Del(ctx, key)
		}
	}
	return nil
}
