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
	rdb.FlushDB(ctx)
	return &RedisCache{
		Client:  rdb,
		MaxSize: maxSize,
	}
}

// REDIS LRU OPERATION METHODS
func (rc *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if ttl == 0 {
		return errors.New("ttl cannot be zero")
	}
	err := rc.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	rc.updateAccessOrder(key)    //Maintain LRU Order
	return rc.evictIfNecessary() //perform LRU eviction if more than size
}

func (rc *RedisCache) Get(key string) (string, error) {
	val, err := rc.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	rc.updateAccessOrder(key)
	return val, nil
}

func (rc *RedisCache) updateAccessOrder(key string) {
	// remove and readd to maintain LRU order
	rc.Client.LRem(ctx, "cache_keys", 0, key)
	rc.Client.LPush(ctx, "cache_keys", key)
}
func (rc *RedisCache) evictIfNecessary() error {
	size := rc.Client.LLen(ctx, "cache_keys").Val() //Redis List length, val return int64
	if size > int64(rc.MaxSize) {                   //
		excess := size - int64(rc.MaxSize)
		for i := int64(0); i < excess; i++ {
			key := rc.Client.RPop(ctx, "cache_keys").Val()
			rc.Client.Del(ctx, key)
		}
	}
	return nil
}

func (rc *RedisCache) GetAll() (map[string]interface{}, error) {
	keys, err := rc.Client.LRange(ctx, "cache_keys", 0, -1).Result() //get all elements from list
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	for _, key := range keys {
		val, err := rc.Client.Get(ctx, key).Result() //returns the value for the specific key
		if err != nil {
			continue
		}
		values[key] = val
	}
	return values, nil
}

func (rc *RedisCache) Delete(key string) error {
	err := rc.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	rc.Client.LRem(ctx, "cache_keys", 0, key) //Remove element from list
	return nil
}
func (rc *RedisCache) DeleteAll() error {
	_, err := rc.Client.FlushAll(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
