package test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/pragadeesh-mcw/Go-Mini-Project/redis_cache"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func setupTestCache() *redis_cache.RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	client.FlushDB(ctx)
	return redis_cache.NewCache("localhost:6379", "", 0, 5)
}
func TestRedisLRUEviction(t *testing.T) {
	cache := setupTestCache()
	//LRU eviction with step wise check
	for i := 0; i < 7; i++ { //greater size to check eviction
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	expectedKeys := []string{"key6", "key5", "key4", "key3", "key2"}
	actualKeys, err := cache.Client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get keys: %v", err)
	}

	if len(actualKeys) != 5 {
		t.Fatalf("Expected 5 keys, got %d", len(actualKeys))
	}

	for i, key := range expectedKeys {
		if actualKeys[i] != key {
			t.Errorf("Expected key %s at position %d, got %s", key, i, actualKeys[i])
		}
	}

	for _, key := range []string{"key0", "key1"} {
		_, err := cache.Get(key)
		if err == nil {
			t.Errorf("Expected key %s to be evicted, but it is still present", key)
		}
	}
}

func TestRedisLRUUpdateAccessOrder(t *testing.T) {
	cache := setupTestCache()
	//redis list order check
	for i := 0; i < 5; i++ {
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	cache.Get("key0")
	cache.Set("key5", "value5", 10*time.Second)

	expectedKeys := []string{"key5", "key0", "key4", "key3", "key2"}
	actualKeys, err := cache.Client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get keys: %v", err)
	}

	if len(actualKeys) != 5 {
		t.Fatalf("Expected 5 keys, got %d", len(actualKeys))
	}

	for i, key := range expectedKeys {
		if actualKeys[i] != key {
			t.Errorf("Expected key %s at position %d, got %s", key, i, actualKeys[i])
		}
	}

	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected key 'key1' to be evicted, but it is still present")
	}
}
func TestRedisSetGet(t *testing.T) {
	cache := setupTestCache()

	err := cache.Set("key1", "value1", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	val, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if val != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}
}
func TestRedisDelete(t *testing.T) {
	cache := setupTestCache()

	err := cache.Set("key1", "value1", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected error for deleted key, got nil")
	}
}

func TestRedisGetAll(t *testing.T) {
	cache := setupTestCache()

	for i := 0; i < 5; i++ {
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	values, err := cache.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all keys: %v", err)
	}

	if len(values) != 5 {
		t.Errorf("Expected 5 keys, got %d", len(values))
	}
}

func TestRedisDeleteAll(t *testing.T) {
	cache := setupTestCache()

	for i := 0; i < 5; i++ {
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	err := cache.DeleteAll()
	if err != nil {
		t.Fatalf("Failed to delete all keys: %v", err)
	}

	values, err := cache.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all keys: %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(values))
	}
}
