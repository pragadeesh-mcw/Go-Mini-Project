package test

import (
	"strconv"
	"testing"
	"time"
	inmemory "unified/in_memory"
)

func BenchmarkInMemory_Set(b *testing.B) {
	cache := inmemory.NewLRUCache(1000, 5)
	for n := 0; n < b.N; n++ {
		cache.Set("key"+strconv.Itoa(n), "value"+strconv.Itoa(n), 10*time.Second)
	}
}

func BenchmarkInMemory_Get(b *testing.B) {
	cache := inmemory.NewLRUCache(1000, 5)
	for n := 0; n < 1000; n++ {
		cache.Set("key"+strconv.Itoa(n), "value"+strconv.Itoa(n), 10*time.Second)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cache.Get("key" + strconv.Itoa(n%1000))
	}
}

func BenchmarkInMemory_Delete(b *testing.B) {
	cache := inmemory.NewLRUCache(1000, 5)
	for n := 0; n < 1000; n++ {
		cache.Set("key"+strconv.Itoa(n), "value"+strconv.Itoa(n), 10*time.Second)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cache.Delete("key" + strconv.Itoa(n%1000))
	}
}

func BenchmarkInMemory_GetAll(b *testing.B) {
	cache := inmemory.NewLRUCache(1000, 5)
	for n := 0; n < 1000; n++ {
		cache.Set("key"+strconv.Itoa(n), "value"+strconv.Itoa(n), 10*time.Second)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cache.GetAll()
	}
}

func BenchmarkInMemory_DeleteAll(b *testing.B) {
	cache := inmemory.NewLRUCache(1000, 5)
	for n := 0; n < 1000; n++ {
		cache.Set("key"+strconv.Itoa(n), "value"+strconv.Itoa(n), 10*time.Second)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cache.DeleteAll()
	}
}
