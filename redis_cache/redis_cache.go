package redis_cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Empty Context
var ctx = context.Background()

// Configurable maxsize and redis.Client Initialization
type RedisCache struct {
	Client  *redis.Client
	MaxSize int
}

// Redis Cache Initialization
func NewCache(addr string, password string, db int, maxSize int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		Client:  rdb,
		MaxSize: maxSize,
	}
}

// REDIS LRU OPERATION METHODS
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if ttl == 0 {
		return errors.New("ttl cannot be zero")
	}
	err := c.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	c.updateAccessOrder(key)    //Maintain LRU Order
	return c.evictIfNecessary() //perform LRU eviction if more than size
}

func (c *RedisCache) Get(key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	c.updateAccessOrder(key)
	return val, nil
}

func (c *RedisCache) updateAccessOrder(key string) {
	// remove and readd to maintain LRU order
	c.Client.LRem(ctx, "cache_keys", 0, key)
	c.Client.LPush(ctx, "cache_keys", key)
}
func (c *RedisCache) evictIfNecessary() error {
	size := c.Client.LLen(ctx, "cache_keys").Val() //Redis List length, val return int64
	if size > int64(c.MaxSize) {                   //
		excess := size - int64(c.MaxSize)
		for i := int64(0); i < excess; i++ {
			key := c.Client.RPop(ctx, "cache_keys").Val()
			c.Client.Del(ctx, key)
		}
	}
	return nil
}

func (c *RedisCache) GetAll() (map[string]string, error) {
	keys, err := c.Client.LRange(ctx, "cache_keys", 0, -1).Result() //get all elements from list
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	for _, key := range keys {
		val, err := c.Client.Get(ctx, key).Result() //returns the value for the specific key
		if err != nil {
			continue
		}
		values[key] = val
	}
	return values, nil
}

func (c *RedisCache) Delete(key string) error {
	err := c.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	c.Client.LRem(ctx, "cache_keys", 0, key) //Remove element from list
	return nil
}

func (c *RedisCache) DeleteAll() error {
	keys, err := c.Client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		c.Client.Del(ctx, key)
	}

	c.Client.Del(ctx, "cache_keys") //delete entire list
	return nil
}
