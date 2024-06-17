package in_memory

import (
	"sync"
	"testing"
)

func BenchmarkLRUCacheSet(b *testing.B) {
	cache := NewLRUCache(1000, 60)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + string(rune(i))
		cache.Set(key, i, 60)
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	cache := NewLRUCache(1000, 60)
	for i := 0; i < 1000; i++ {
		key := "key" + string(rune(i))
		cache.Set(key, i, 60)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + string(rune(i%1000))
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
			key := "key" + string(rune(idx))
			cache.Set(key, idx, 60)
		}(i)
	}
	wg.Wait()
}

func BenchmarkLRUCacheConcurrentGet(b *testing.B) {
	cache := NewLRUCache(1000, 60)
	for i := 0; i < 1000; i++ {
		key := "key" + string(rune(i))
		cache.Set(key, i, 60)
	}

	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key" + string(rune(idx%1000))
			cache.Get(key)
		}(i)
	}
	wg.Wait()
}
