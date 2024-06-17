package redis_test

import (
	"strconv"
	"testing"
	"time"

	redis "unified/redis"
)

var (
	addr     = "localhost:6379"
	password = ""
	db       = 0
)

func setup() {
	redis.InitRedis(addr, password, db)
}

func BenchmarkSet(b *testing.B) {
	setup()
	cache := redis.NewCache(addr, password, db)

	b.ResetTimer() //to exclude the setup overhead
	for i := 0; i < b.N; i++ {
		cache.Set("bench-key", "bench-value", 10*time.Second)
	}
}

func BenchmarkGet(b *testing.B) {
	setup()
	cache := redis.NewCache(addr, password, db)
	cache.Set("bench-key", "bench-value", 10*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("bench-key")
	}
}

func BenchmarkGetAll(b *testing.B) {
	setup()
	cache := redis.NewCache(addr, password, db)
	for i := 0; i < 100; i++ {
		cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAll()
	}
}

func BenchmarkDelete(b *testing.B) {
	setup()
	cache := redis.NewCache(addr, password, db)
	cache.Set("bench-key", "bench-value", 10*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete("bench-key")
		cache.Set("bench-key", "bench-value", 10*time.Second)
	}
}

func BenchmarkDeleteAll(b *testing.B) {
	setup()
	cache := redis.NewCache(addr, password, db)
	for i := 0; i < 100; i++ {
		cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 10*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.DeleteAll()
		for j := 0; j < 100; j++ {
			cache.Set("key"+strconv.Itoa(j), "value"+strconv.Itoa(j), 10*time.Second)
		}
	}
}
