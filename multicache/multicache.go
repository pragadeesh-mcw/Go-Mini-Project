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

func (c *MultiCache) Set(key string, value interface{}, ttl time.Duration) error {
	// Set in both caches
	c.inMemoryCache.Set(key, value, ttl)
	err := c.redisCache.Set(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (c *MultiCache) Get(key string) (interface{}, error) {
	value1, _ := c.inMemoryCache.Get(key)
	value2, err := c.redisCache.Get(key)

	if value1 == value2 {
		return value2, err
	}
	if err != nil {
		return nil, err
	}
	return value2, nil
}

func (c *MultiCache) GetAll() (map[string]interface{}, error) {
	// Get all from Redis
	redisValues, err := c.redisCache.GetAll()
	if err != nil {
		return nil, err
	}

	// Get all from in-memory cache
	inMemoryValues := c.inMemoryCache.GetAll()

	// Check if redisValues and inMemoryValues are the same
	if reflect.DeepEqual(redisValues, inMemoryValues) {
		return inMemoryValues, nil
	}
	return inMemoryValues, nil
}

func (c *MultiCache) Delete(key string) error {
	// Delete from both caches
	c.inMemoryCache.Delete(key)
	err := c.redisCache.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (c *MultiCache) DeleteAll() error {
	// Delete all from both caches
	c.inMemoryCache.DeleteAll()
	err := c.redisCache.DeleteAll()
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
