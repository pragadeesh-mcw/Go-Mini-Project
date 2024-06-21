package test

import (
	"context"
	"strconv"
	"testing"
	"time"
	"unified/redis_cache"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func setupRedisTestCache() *redis_cache.RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	client.FlushDB(ctx)
	return redis_cache.NewCache("localhost:6379", "", 0, 5)
}

// 1. Test LRU Eviction
func TestRedisLRUEviction(t *testing.T) {
	cache := setupRedisTestCache()
	// LRU eviction with step-wise check
	for i := 0; i < 7; i++ { // greater size to check eviction
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

// 2. Test LRU Access Order
func TestRedisLRUUpdateAccessOrder(t *testing.T) {
	cache := setupRedisTestCache()
	// Redis list order check
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

// 3. Test Set and Get
func TestRedisSetGet(t *testing.T) {
	cache := setupRedisTestCache()

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

// 4. Test Delete
func TestRedisDelete(t *testing.T) {
	cache := setupRedisTestCache()

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

// 5. Test Getall
func TestRedisGetAll(t *testing.T) {
	cache := setupRedisTestCache()

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

// 6. Test Deleteall
func TestRedisDeleteAll(t *testing.T) {
	cache := setupRedisTestCache()

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

// Edge and Corner Test Cases
// 7. Test Set Empty Key
func TestRedisEmptyKey(t *testing.T) {
	cache := setupRedisTestCache()

	err := cache.Set("", "value", 10*time.Second)
	if err == nil {
		t.Error("Expected error when setting an empty key, got nil")
	}

	_, err = cache.Get("")
	if err == nil {
		t.Error("Expected error when getting an empty key, got nil")
	}
}

// 8. Test Set and Get Large value
func TestRedisLargeValue(t *testing.T) {
	cache := setupRedisTestCache()
	largeValue := make([]byte, 1024*1024) // 1 MB value

	err := cache.Set("key_large", largeValue, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set large value: %v", err)
	}

	val, err := cache.Get("key_large")
	if err != nil {
		t.Fatalf("Failed to get large value: %v", err)
	}

	if string(val) != string(largeValue) {
		t.Errorf("Expected large value, got something else")
	}
}

// 9. Test Immediate Expiration
func TestRedisImmediateExpiration(t *testing.T) {
	cache := setupRedisTestCache()

	err := cache.Set("key1", "value1", 0)
	if err == nil {
		t.Errorf("Expected error when setting with immediate expiration, got nil")
	}

	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected error when getting immediately expired key, got nil")
	}
}

// 10. Test Short TTL
func TestRedisShortTTL(t *testing.T) {
	cache := setupRedisTestCache()

	err := cache.Set("key1", "value1", 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to set key with short TTL: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected key1 to expire, but it is still present")
	}
}

// 11. Test Frequent Access
func TestRedisLRUWithFrequentAccess(t *testing.T) {
	cache := setupRedisTestCache()

	// Set initial keys
	for i := 0; i < 5; i++ {
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	// Frequently access key0 to prevent eviction
	for i := 0; i < 10; i++ {
		_, err := cache.Get("key0")
		if err != nil {
			t.Fatalf("Failed to get key0: %v", err)
		}
	}

	// Add a new key to trigger eviction
	err := cache.Set("key5", "value5", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	_, err = cache.Get("key0")
	if err != nil {
		t.Errorf("Expected key0 to be present due to frequent access, but it was evicted")
	}

	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected key1 to be evicted, but it is still present")
	}
}

// 12. Test Simultaneous Set and Get
func TestRedisLRUWithSimultaneousSetAndGet(t *testing.T) {
	cache := setupRedisTestCache()

	// Set initial keys
	for i := 0; i < 5; i++ {
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	done := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 10*time.Second)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 5; i++ {
			cache.Get("key" + strconv.Itoa(i))
		}
		done <- true
	}()

	<-done
	<-done

	// Verify the keys after concurrent set and get operations
	for i := 0; i < 5; i++ {
		_, err := cache.Get("key" + strconv.Itoa(i))
		if err != nil {
			t.Errorf("Expected key%d to be present, but it was evicted", i)
		}
	}
}

// 13. Test Max Capacity
func TestRedisSetWithMaxCapacity(t *testing.T) {
	cache := setupRedisTestCache()

	for i := 0; i < 5; i++ {
		key := "key" + strconv.Itoa(i)
		err := cache.Set(key, "value"+strconv.Itoa(i), 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	}

	// Try to set an additional key beyond the capacity
	err := cache.Set("key5", "value5", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	keys, err := cache.Client.LRange(ctx, "cache_keys", 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get keys: %v", err)
	}

	if len(keys) != 5 {
		t.Fatalf("Expected 5 keys, got %d", len(keys))
	}
}

// 14. Test Set Same Key
func TestRedisSetWithSameKey(t *testing.T) {
	cache := setupRedisTestCache()

	err := cache.Set("key1", "value1", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	err = cache.Set("key1", "new_value1", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	val, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if val != "new_value1" {
		t.Errorf("Expected 'new_value1', got '%s'", val)
	}
}
