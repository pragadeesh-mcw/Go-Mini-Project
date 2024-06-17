package redis_test

import (
	"testing"
	"time"
	redis "unified/redis"
)

func setupTest() {
	redis.InitRedis(addr, password, db)
}

func TestSet(t *testing.T) {
	setupTest()
	cache := redis.NewCache(addr, password, db)
	cache.Set("test-key", "test-value", 10*time.Second)

	val, exists := cache.Get("test-key")
	if !exists {
		t.Errorf("Expected key to exist")
	}
	if val != "test-value" {
		t.Errorf("Expected value to be 'test-value', got '%s'", val)
	}
}

func TestGet(t *testing.T) {
	setupTest()
	cache := redis.NewCache(addr, password, db)
	cache.Set("test-key", "test-value", 10*time.Second)

	val, exists := cache.Get("test-key")
	if !exists {
		t.Errorf("Expected key to exist")
	}
	if val != "test-value" {
		t.Errorf("Expected value to be 'test-value', got '%s'", val)
	}

	val, exists = cache.Get("non-existent-key")
	if exists {
		t.Errorf("Expected key to not exist")
	}
	if val != "" {
		t.Errorf("Expected value to be an empty string, got '%s'", val)
	}
}

// func TestGetAll(t *testing.T) {
// 	setupTest()
// 	cache := redis.NewCache(addr, password, db)
// 	cache.Set("key1", "value1", 10*time.Second)
// 	cache.Set("key2", "value2", 10*time.Second)

// 	all := cache.GetAll()
// 	if len(all) != 2 {
// 		t.Errorf("Expected 2 keys in the cache, got %d", len(all))
// 	}
// 	if all["key1"] != "value1" {
// 		t.Errorf("Expected value to be 'value1', got '%s'", all["key1"])
// 	}
// 	if all["key2"] != "value2" {
// 		t.Errorf("Expected value to be 'value2', got '%s'", all["key2"])
// 	}
// }

func TestLRU(t *testing.T) {
	setupTest()
	cache := redis.NewCache(addr, password, db)

	cache.Set("key1", "value1", 10*time.Second)
	cache.Set("key2", "value2", 10*time.Second)
	cache.Set("key3", "value3", 10*time.Second)
	cache.Get("key1")
	cache.Get("key2")

	cache.Set("key4", "value4", 10*time.Second)

	_, exists := cache.Get("key3")
	if exists {
		t.Errorf("Expected key3 to be evicted due to LRU, but it still exists in cache")
	}

	keys := []string{"key1", "key2", "key4"}
	for _, key := range keys {
		_, exists := cache.Get(key)
		if !exists {
			t.Errorf("Expected key %s to exist in cache, but it doesn't", key)
		}
	}
}

func TestDelete(t *testing.T) {
	setupTest()
	cache := redis.NewCache(addr, password, db)
	cache.Set("test-key", "test-value", 10*time.Second)

	cache.Delete("test-key")
	val, exists := cache.Get("test-key")
	if exists {
		t.Errorf("Expected key to not exist")
	}
	if val != "" {
		t.Errorf("Expected value to be an empty string, got '%s'", val)
	}
}

func TestDeleteAll(t *testing.T) {
	setupTest()
	cache := redis.NewCache(addr, password, db)
	cache.Set("key1", "value1", 10*time.Second)
	cache.Set("key2", "value2", 10*time.Second)

	cache.DeleteAll()
	all := cache.GetAll()
	if len(all) != 0 {
		t.Errorf("Expected 0 keys in the cache, got %d", len(all))
	}
}
