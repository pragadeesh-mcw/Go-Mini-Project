package multicache

import (
	"reflect"
	"time"

	"unified/in_memory"
	"unified/redis_cache"
)

type MultiCache struct {
	inMemoryCache *in_memory.LRUCache
	redisCache    *redis_cache.RedisCache
}

func NewMultiCache(inMemoryCache *in_memory.LRUCache, redisCache *redis_cache.RedisCache) *MultiCache {
	return &MultiCache{
		inMemoryCache: inMemoryCache,
		redisCache:    redisCache,
	}
}

func (mc *MultiCache) Set(key string, value interface{}, ttl time.Duration) error {
	// Set in both caches
	err := mc.redisCache.Set(key, value, ttl)
	mc.inMemoryCache.Set(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MultiCache) Get(key string) (interface{}, error) {
	value1, _ := mc.inMemoryCache.Get(key)
	value2, err := mc.redisCache.Get(key)

	if value1 == value2 {
		return value2, err
	}
	if err != nil {
		return nil, err
	}
	return value2, nil
}

func (mc *MultiCache) GetAll() (map[string]interface{}, error) {
	// Get all from Redis
	redisValues, err := mc.redisCache.GetAll()
	if err != nil {
		return nil, err
	}

	// Get all from in-memory cache
	inMemoryValues := mc.inMemoryCache.GetAll()

	// Check if redisValues and inMemoryValues are the same
	if reflect.DeepEqual(redisValues, inMemoryValues) {
		return redisValues, nil
	}
	return redisValues, nil
}

func (mc *MultiCache) Delete(key string) error {
	// Delete from both caches
	mc.inMemoryCache.Delete(key)
	err := mc.redisCache.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MultiCache) DeleteAll() error {
	// Delete all from both caches
	mc.inMemoryCache.DeleteAll()
	err := mc.redisCache.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}

// func (c *MultiCache) EvictFromBothCaches(key string) {
// 	// Evict key from both caches
// 	c.inMemoryCache.Delete(key)
// 	c.redisCache.Delete(key)
// }
