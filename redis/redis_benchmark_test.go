package redis

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var benchCtx = context.Background()

func setupBenchmarkCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	client.FlushDB(benchCtx)
	return NewCache("localhost:6379", "", 0, 1000)
}

func BenchmarkCacheSet(b *testing.B) {
	cache := setupBenchmarkCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys.
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%1000) // Retrieve one of the preloaded keys.
		cache.Get(key)
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys.
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%1000) // Delete one of the preloaded keys.
		cache.Delete(key)
	}
}

func BenchmarkCacheGetAll(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys.
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAll()
	}
}

func BenchmarkCacheLRUEviction(b *testing.B) {
	cache := setupBenchmarkCache()
	cache.maxSize = 100 // Set a small maxSize to trigger evictions during benchmark.

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}
}
