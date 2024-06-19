package test

import (
	"testing"
	"time"

	"unified/in_memory"
	"unified/multicache"
	"unified/redis_cache"
)

func TestMultiCache(t *testing.T) {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, 100)

	cache := multicache.NewMultiCache(inMemoryCache, redisCache)

	key := "test-key"
	value := "test-value"
	ttl := 10 * time.Second

	// Test Set
	err := cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Test Get from in-memory cache
	cachedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get key from cache: %v", err)
	}
	if cachedValue != value {
		t.Fatalf("Expected value %v, got %v", value, cachedValue)
	}

	// Test Get from Redis cache, after evicting from in-memory
	inMemoryCache.Delete(key)
	cachedValue, err = cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get key from Redis cache: %v", err)
	}
	if cachedValue != value {
		t.Fatalf("Expected value %v, got %v", value, cachedValue)
	}

	// Test GetAll
	allValues, err := cache.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all values: %v", err)
	}
	if len(allValues) == 0 {
		t.Fatalf("Expected non-zero length of all values")
	}

	// Test Delete
	err = cache.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}
	_, err = cache.Get(key)
	if err == nil {
		t.Fatalf("Expected error when getting deleted key")
	}

	// Test DeleteAll
	err = cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}
	err = cache.DeleteAll()
	if err != nil {
		t.Fatalf("Failed to delete all keys: %v", err)
	}
	_, err = cache.Get(key)
	if err == nil {
		t.Fatalf("Expected error when getting key after deleting all keys")
	}
}
