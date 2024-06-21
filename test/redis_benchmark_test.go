package test

import (
	"context"
	"strconv"
	"testing"
	"time"
	"unified/redis_cache"

	"github.com/redis/go-redis/v9"
)

var benchCtx = context.Background()

func setupBenchmarkCache() *redis_cache.RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	client.FlushDB(benchCtx) //start with empty redis
	return redis_cache.NewCache("localhost:6379", "", 0, 1000)
}

func BenchmarkRedisSet(b *testing.B) {
	cache := setupBenchmarkCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}
}

func BenchmarkRedisGet(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%1000) // Retrieve one of the preloaded keys
		cache.Get(key)
	}
}
func BenchmarkRedisGetAll(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAll() //Retrieve all keys
	}
}
func BenchmarkRedisDelete(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%1000) // Delete one of the preloaded keys
		cache.Delete(key)
	}
}

func BenchmarkRedisDeleteAll(b *testing.B) {
	cache := setupBenchmarkCache()

	// Preload the cache with keys
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := cache.DeleteAll()
		if err != nil {
			b.Fatalf("Failed to delete all keys: %v", err)
		}

		// Re-populate the cache with keys for the next iteration
		for i := 0; i < 1000; i++ {
			key := "key" + strconv.Itoa(i)
			cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		}
	}
}
