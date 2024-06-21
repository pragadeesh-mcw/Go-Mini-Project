package test

import (
	"strconv"
	"testing"
	"time"
	inmemory "unified/in_memory"
)

func TestInMemoryCache(t *testing.T) {
	cache := inmemory.NewLRUCache(3, 60)

	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	cache.Set("key3", "value3", 30)
	//1. GET EXISTING
	t.Run("Get existing key", func(t *testing.T) {
		value, ok := cache.Get("key1")
		if !ok {
			t.Errorf("Expected key1 to exist in cache")
		}
		if value != "value1" {
			t.Errorf("Expected value 'value1', got %v", value)
		}
	})
	//2. GET NON-EXISTING
	t.Run("Get non-existing key", func(t *testing.T) {
		_, ok := cache.Get("key4")
		if ok {
			t.Errorf("Expected key4 not to exist in cache")
		}
	})
	//3. GET ALL
	t.Run("GetAll", func(t *testing.T) {
		items := cache.GetAll()
		expected := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for key, expectedValue := range expected {
			if items[key] != expectedValue {
				t.Errorf("Expected value '%v' for key '%s', got %v", expectedValue, key, items[key])
			}
		}
	})
	//4. DELETE KEY
	t.Run("Delete key", func(t *testing.T) {
		cache.Delete("key2")
		_, ok := cache.Get("key2")
		if ok {
			t.Errorf("Expected key2 to be deleted from cache")
		}
	})
	//5. DELETE ALL KEYS
	t.Run("DeleteAll", func(t *testing.T) {
		cache.DeleteAll()
		items := cache.GetAll()
		if len(items) != 0 {
			t.Errorf("Expected cache to be empty after DeleteAll, got %v items", len(items))
		}
	})

	cache = inmemory.NewLRUCache(3, 1) // Setting a small TTL for quick eviction tests
	//6. FASTER GET AND SET
	t.Run("Set and Get with TTL expiration", func(t *testing.T) {
		cache.Set("key1", "value1", 1)
		time.Sleep(2 * time.Second) // wait for the key to expire
		_, ok := cache.Get("key1")
		if ok {
			t.Errorf("Expected key1 to expire and not exist in cache")
		}
	})
	//7. LRU ON CAPACITY OVERFLOW
	t.Run("LRU eviction on capacity overflow", func(t *testing.T) {
		cache.Set("key1", "value1", 10)
		cache.Set("key2", "value2", 10)
		cache.Set("key3", "value3", 10)
		cache.Set("key4", "value4", 10) // this should evict "key1" as it's LRU
		_, ok := cache.Get("key1")
		if ok {
			t.Errorf("Expected key1 to be evicted due to capacity overflow")
		}
	})
	//8. UPDATE KEY
	t.Run("Update existing key", func(t *testing.T) {
		cache.Set("key2", "new_value2", 10)
		value, ok := cache.Get("key2")
		if !ok {
			t.Errorf("Expected key2 to exist in cache")
		}
		if value != "new_value2" {
			t.Errorf("Expected value 'new_value2', got %v", value)
		}
	})
	//9. IMMEDIATE EXPIRATION
	t.Run("Set with immediate expiration", func(t *testing.T) {
		cache.Set("key5", "value5", 0)
		_, ok := cache.Get("key5")
		if ok {
			t.Errorf("Expected key5 to expire immediately and not exist in cache")
		}
	})
	//10. GETALL AFTER EXPIRATION
	t.Run("GetAll after some keys expire", func(t *testing.T) {
		cache.Set("key6", "value6", 1)
		cache.Set("key7", "value7", 2)
		time.Sleep(3 * time.Second) // wait for keys to expire
		items := cache.GetAll()
		if len(items) != 0 {
			t.Errorf("Expected no items in cache, got %v items", len(items))
		}
	})
	//11. DELETE NON-EXISTING
	t.Run("Delete non-existing key", func(t *testing.T) {
		deleted := cache.Delete("non_existing_key")
		if deleted {
			t.Errorf("Expected delete to return false for non-existing key")
		}
	})
	//12. EVICT AND ACCESS
	t.Run("Evict when accessing existing key", func(t *testing.T) {
		cache.Set("key8", "value8", 10)
		cache.Set("key9", "value9", 10)
		cache.Set("key10", "value10", 10)
		cache.Get("key8") // Access "key8" to make it recently used
		cache.Set("key11", "value11", 10)
		// Should evict "key9" now
		_, ok := cache.Get("key9")
		if ok {
			t.Errorf("Expected key9 to be evicted due to LRU policy")
		}
	})
	//13. DELETE AND RE-ADD
	t.Run("Delete key and re-add", func(t *testing.T) {
		cache.Set("key11", "value11", 10)
		cache.Delete("key11")
		cache.Set("key11", "new_value11", 10)
		value, ok := cache.Get("key11")
		if !ok {
			t.Errorf("Expected key11 to exist in cache")
		}
		if value != "new_value11" {
			t.Errorf("Expected value 'new_value11', got %v", value)
		}
	})
	//14. DELETE WHEN EMPTY
	t.Run("DeleteAll with empty cache", func(t *testing.T) {
		cache.DeleteAll()
		cache.DeleteAll()
		items := cache.GetAll()
		if len(items) != 0 {
			t.Errorf("Expected cache to be empty, got %v items", len(items))
		}
	})
	//15. CONCURRENT SET AND GET
	t.Run("Concurrency test", func(t *testing.T) {
		cache := inmemory.NewLRUCache(5, 1)
		done := make(chan bool)

		go func() {
			for i := 0; i < 100; i++ {
				cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 10)
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				cache.Get("key" + strconv.Itoa(i))
			}
			done <- true
		}()

		<-done
		<-done

		if len(cache.GetAll()) > 5 {
			t.Errorf("Expected at most 5 items in cache")
		}
	})
}
