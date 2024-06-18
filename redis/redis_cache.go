package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Empty Context
var ctx = context.Background()

// Configurable maxsize and redis.Client Initialization
type RedisCache struct {
	client  *redis.Client
	maxSize int
}

// Redis Cache Initialization
func NewCache(addr string, password string, db int, maxSize int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client:  rdb,
		maxSize: maxSize,
	}
}

// REDIS LRU OPERATION METHODS
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	c.updateAccessOrder(key)    //Maintain LRU Order
	return c.evictIfNecessary() //perform LRU eviction if more than size
}

func (c *RedisCache) Get(key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	c.updateAccessOrder(key)
	return val, nil
}

func (c *RedisCache) updateAccessOrder(key string) {
	// remove and readd to maintain LRU order
	c.client.LRem(ctx, "cache_keys", 0, key)
	c.client.LPush(ctx, "cache_keys", key)
}
func (c *RedisCache) evictIfNecessary() error {
	size := c.client.LLen(ctx, "cache_keys").Val() //Redis List length
	if size > int64(c.maxSize) {
		excess := size - int64(c.maxSize)
		for i := int64(0); i < excess; i++ {
			key := c.client.RPop(ctx, "cache_keys").Val()
			c.client.Del(ctx, key)
		}
	}
	return nil
}

func (c *RedisCache) GetAll() (map[string]string, error) {
	keys, err := c.client.LRange(ctx, "cache_keys", 0, -1).Result() //get all elements from list
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

func (c *RedisCache) Delete(key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	c.client.LRem(ctx, "cache_keys", 0, key) //Remove element from list
	return nil
}

func (c *RedisCache) DeleteAll() error {
	keys, err := c.client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		c.client.Del(ctx, key)
	}

	c.client.Del(ctx, "cache_keys") //delete entire list
	return nil
}
