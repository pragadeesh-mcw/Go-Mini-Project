package multicache

import (
	"strconv"
	"testing"
	"time"

	"unified/in_memory"
	"unified/redis"
)

func BenchmarkMultiCacheSet(b *testing.B) {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis.NewCache("localhost:6379", "", 0, 100)

	cache := NewMultiCache(inMemoryCache, redisCache)

	key := "test-key"
	value := "test-value"
	ttl := 10 * time.Second

	for i := 0; i < b.N; i++ {
		err := cache.Set(key, value, ttl)
		if err != nil {
			b.Fatalf("Failed to set key: %v", err)
		}
	}
}

func BenchmarkMultiCacheGet(b *testing.B) {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis.NewCache("localhost:6379", "", 0, 100)

	cache := NewMultiCache(inMemoryCache, redisCache)

	key := "test-key"
	value := "test-value"
	ttl := 10 * time.Second

	err := cache.Set(key, value, ttl)
	if err != nil {
		b.Fatalf("Failed to set key: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := cache.Get(key)
		if err != nil {
			b.Fatalf("Failed to get key: %v", err)
		}
	}
}

func BenchmarkMultiCacheGetAll(b *testing.B) {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis.NewCache("localhost:6379", "", 0, 100)

	cache := NewMultiCache(inMemoryCache, redisCache)

	for i := 0; i < 100; i++ {
		key := "test-key" + strconv.Itoa(i)
		value := "test-value"
		ttl := 10 * time.Second

		err := cache.Set(key, value, ttl)
		if err != nil {
			b.Fatalf("Failed to set key: %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := cache.GetAll()
		if err != nil {
			b.Fatalf("Failed to get all values: %v", err)
		}
	}
}
func BenchmarkMultiCacheDelete(b *testing.B) {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis.NewCache("localhost:6379", "", 0, 100)

	cache := NewMultiCache(inMemoryCache, redisCache)

	key := "test-key"
	value := "test-value"
	ttl := 10 * time.Second

	err := cache.Set(key, value, ttl)
	if err != nil {
		b.Fatalf("Failed to set key: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := cache.Delete(key)
		if err != nil {
			b.Fatalf("Failed to delete key: %v", err)
		}
		cache.Set(key, value, ttl) // Re-set the key for the next iteration
	}
}

func BenchmarkMultiCacheDeleteAll(b *testing.B) {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis.NewCache("localhost:6379", "", 0, 100)

	cache := NewMultiCache(inMemoryCache, redisCache)

	for i := 0; i < 100; i++ {
		key := "test-key" + strconv.Itoa(i)
		value := "test-value"
		ttl := 10 * time.Second

		err := cache.Set(key, value, ttl)
		if err != nil {
			b.Fatalf("Failed to set key: %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := cache.DeleteAll()
		if err != nil {
			b.Fatalf("Failed to delete all keys: %v", err)
		}

		// Re-set the keys for the next iteration
		for i := 0; i < 100; i++ {
			key := "test-key" + strconv.Itoa(i)
			value := "test-value"
			ttl := 10 * time.Second

			err := cache.Set(key, value, ttl)
			if err != nil {
				b.Fatalf("Failed to set key: %v", err)
			}
		}
	}
}
