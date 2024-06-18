package in_memory

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(3, 60)

	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	cache.Set("key3", "value3", 30)

	t.Run("Get existing key", func(t *testing.T) {
		value, ok := cache.Get("key1")
		if !ok {
			t.Errorf("Expected key1 to exist in cache")
		}
		if value != "value1" {
			t.Errorf("Expected value 'value1', got %v", value)
		}
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		_, ok := cache.Get("key4")
		if ok {
			t.Errorf("Expected key4 not to exist in cache")
		}
	})

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

	t.Run("Delete key", func(t *testing.T) {
		cache.Delete("key2")
		_, ok := cache.Get("key2")
		if ok {
			t.Errorf("Expected key2 to be deleted from cache")
		}
	})

	t.Run("DeleteAll", func(t *testing.T) {
		cache.DeleteAll()
		items := cache.GetAll()
		if len(items) != 0 {
			t.Errorf("Expected cache to be empty after DeleteAll, got %v items", len(items))
		}
	})
}
