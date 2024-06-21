package test

import (
	"strconv"
	"testing"
	"time"

	"unified/in_memory"
	"unified/multicache"
	"unified/redis_cache"
)

func setupTestInMemoryCache() *in_memory.LRUCache {
	return in_memory.NewLRUCache(10, 60)
}

func setupTestRedisCache() *redis_cache.RedisCache {
	return redis_cache.NewCache("localhost:6379", "", 0, 10)
}

func TestMultiCache_SetGet(t *testing.T) {
	inMemoryCache := setupTestInMemoryCache()
	redisCache := setupTestRedisCache()
	cache := multicache.NewMultiCache(inMemoryCache, redisCache)

	key := "test-key"
	value := "test-value"
	ttl := 10 * time.Second

	err := cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Get from in-memory cache
	cachedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get key from cache: %v", err)
	}
	if cachedValue != value {
		t.Fatalf("Expected value %v, got %v", value, cachedValue)
	}

	// Get from Redis cache after eviction from in-memory cache
	inMemoryCache.Delete(key)
	cachedValue, err = cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get key from Redis cache: %v", err)
	}
	if cachedValue != value {
		t.Fatalf("Expected value %v, got %v", value, cachedValue)
	}
}

func TestMultiCache_GetAll(t *testing.T) {
	inMemoryCache := setupTestInMemoryCache()
	redisCache := setupTestRedisCache()
	cache := multicache.NewMultiCache(inMemoryCache, redisCache)

	// Set some keys in the cache
	for i := 1; i <= 5; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		err := cache.Set(key, value, 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key %s: %v", key, err)
		}
	}

	// Get all keys from MultiCache
	allValues, err := cache.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all values: %v", err)
	}

	// Check if all values retrieved are correct
	expectedValues := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	if len(allValues) != len(expectedValues) {
		t.Fatalf("Expected %d keys, got %d", len(expectedValues), len(allValues))
	}

	for key, expectedValue := range expectedValues {
		if allValues[key] != expectedValue {
			t.Errorf("Expected value %v for key %s, got %v", expectedValue, key, allValues[key])
		}
	}
}

func TestMultiCache_Delete(t *testing.T) {
	inMemoryCache := setupTestInMemoryCache()
	redisCache := setupTestRedisCache()
	cache := multicache.NewMultiCache(inMemoryCache, redisCache)

	key := "test-key"
	value := "test-value"
	ttl := 10 * time.Second

	// Set a key-value pair
	err := cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Delete the key
	err = cache.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Ensure key is deleted from both caches
	_, err = cache.Get(key)
	if err == nil {
		t.Errorf("Expected error when getting deleted key, got nil")
	}
}

func TestMultiCache_DeleteAll(t *testing.T) {
	inMemoryCache := setupTestInMemoryCache()
	redisCache := setupTestRedisCache()
	cache := multicache.NewMultiCache(inMemoryCache, redisCache)

	// Set some keys in the cache
	for i := 1; i <= 5; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		err := cache.Set(key, value, 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key %s: %v", key, err)
		}
	}

	// Delete all keys
	err := cache.DeleteAll()
	if err != nil {
		t.Fatalf("Failed to delete all keys: %v", err)
	}

	// Ensure all keys are deleted from both caches
	for i := 1; i <= 5; i++ {
		key := "key" + strconv.Itoa(i)
		_, err := cache.Get(key)
		if err == nil {
			t.Errorf("Expected error when getting deleted key %s, got nil", key)
		}
	}
}

func TestMultiCache_TTLExpiration(t *testing.T) {
	inMemoryCache := setupTestInMemoryCache()
	redisCache := setupTestRedisCache()
	cache := multicache.NewMultiCache(inMemoryCache, redisCache)

	key := "ttl-key"
	value := "ttl-value"
	shortTTL := 1 * time.Second
	longTTL := 10 * time.Second

	// Set with short TTL
	err := cache.Set(key, value, shortTTL)
	if err != nil {
		t.Fatalf("Failed to set key with short TTL: %v", err)
	}

	// Wait for short TTL to expire
	time.Sleep(shortTTL + 100*time.Millisecond)

	// Set with long TTL
	err = cache.Set(key, value, longTTL)
	if err != nil {
		t.Fatalf("Failed to set key with long TTL: %v", err)
	}

	// Ensure key is present in both caches with long TTL
	cachedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get key with long TTL: %v", err)
	}
	if cachedValue != value {
		t.Fatalf("Expected value %v, got %v", value, cachedValue)
	}

	// Wait for long TTL to expire in Redis only
	time.Sleep(longTTL + 100*time.Millisecond)

	// Wait for in-memory cache to evict
	time.Sleep(60 + 100*time.Millisecond)
}
