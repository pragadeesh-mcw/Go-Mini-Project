package in_memory

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkLRUCacheSet(b *testing.B) {
	cache := NewLRUCache(1000, 60)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, i, 60)
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	cache := NewLRUCache(1000, 60)
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, i, 60)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%1000)
		cache.Get(key)
	}
}

func BenchmarkLRUCacheConcurrentSet(b *testing.B) {
	cache := NewLRUCache(1000, 60)
	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(idx)
			cache.Set(key, idx, 60)
		}(i)
	}
	wg.Wait()
}

func BenchmarkLRUCacheConcurrentGet(b *testing.B) {
	cache := NewLRUCache(1000, 60)
	for i := 0; i < 1000; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, i, 60)
	}

	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(idx%1000)
			cache.Get(key)
		}(i)
	}
	wg.Wait()
}
