package multicache

import (
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
	// Try to get from in-memory cache first
	if value, found := c.inMemoryCache.Get(key); found {
		return value, nil
	}

	// Fallback to Redis cache
	value, err := c.redisCache.Get(key)
	if err != nil {
		return nil, err
	}

	// Update in-memory cache
	c.inMemoryCache.Set(key, value, time.Duration(0))
	return value, nil
}

func (c *MultiCache) GetAll() (map[string]interface{}, error) {
	// Get all from Redis
	redisValues, err := c.redisCache.GetAll()
	if err != nil {
		return nil, err
	}

	// Update in-memory cache with all values from Redis
	for key, value := range redisValues {
		c.inMemoryCache.Set(key, value, time.Duration(0))
	}

	// Get all from in-memory cache
	inMemoryValues := c.inMemoryCache.GetAll()

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
