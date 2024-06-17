package in_memory

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	key        string
	value      interface{}
	expiration int64
}

type LRUCache struct {
	capacity int
	ttl      int64
	items    map[string]*list.Element
	order    *list.List
	mutex    sync.Mutex
	evictCh  chan string
}

func NewLRUCache(capacity int, ttl int64) *LRUCache {
	c := &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[string]*list.Element),
		order:    list.New(),
		evictCh:  make(chan string, capacity),
	}
	go c.startEvictionRoutine()
	return c
}

func (c *LRUCache) startEvictionRoutine() {
	ticker := time.NewTicker(time.Duration(c.ttl) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.evictExpired()
		case key := <-c.evictCh:
			c.delete(key)
		}
	}
}

func (c *LRUCache) evictExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().Unix()
	for key, el := range c.items {
		if el.Value.(*CacheItem).expiration < now {
			fmt.Printf("Evicting expired key: %s\n", key)
			delete(c.items, key)
			c.order.Remove(el)
		}
	}
}

func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expirationTime := time.Now().Unix() + int64(expiration.Seconds())

	if el, ok := c.items[key]; ok {
		//if exists, update existing
		c.order.MoveToFront(el)
		el.Value.(*CacheItem).value = value
		el.Value.(*CacheItem).expiration = expirationTime
		fmt.Printf("Updated key: %s, value: %v, expiration: %d\n", key, value, expirationTime)
		return
	}

	if c.order.Len() >= c.capacity {
		c.evict()
	}
	//add new item
	item := &CacheItem{
		key:        key,
		value:      value,
		expiration: expirationTime,
	}
	el := c.order.PushFront(item)
	c.items[key] = el
	fmt.Printf("Set key: %s, value: %v, expiration: %d\n", key, value, expirationTime)
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	el, ok := c.items[key]
	if !ok {
		fmt.Printf("Get key: %s not found\n", key)
		return nil, false
	}

	now := time.Now().Unix()
	if el.Value.(*CacheItem).expiration < now {
		fmt.Printf("Get key: %s expired\n", key)
		c.delete(key)
		return nil, false
	}

	c.order.MoveToFront(el)
	fmt.Printf("Get key: %s, value: %v\n", key, el.Value.(*CacheItem).value)
	return el.Value.(*CacheItem).value, true
}

func (c *LRUCache) GetAll() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make(map[string]interface{})
	now := time.Now().Unix()
	for key, el := range c.items {
		if el.Value.(*CacheItem).expiration >= now {
			result[key] = el.Value.(*CacheItem).value
		} else {
			go func(k string) { c.evictCh <- k }(key) //startEvictionRotine
		}
	}
	return result
}

func (c *LRUCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if el, ok := c.items[key]; ok {
		delete(c.items, key)
		c.order.Remove(el)
		fmt.Printf("Deleted key: %s\n", key)
		return true
	} else {
		fmt.Println("Key is not found")
		return false
	}
}

func (c *LRUCache) DeleteAll() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.items) == 0 {
		fmt.Println("No keys found")
		return false
	}

	c.items = make(map[string]*list.Element)
	c.order.Init() //delete list
	fmt.Println("Deleted all keys")
	return true
}

func (c *LRUCache) evict() {
	el := c.order.Back() //access the LRU element
	if el != nil {
		c.order.Remove(el)
		item := el.Value.(*CacheItem)
		delete(c.items, item.key)
		fmt.Printf("Evicted key: %s\n", item.key)
	}
}

func (c *LRUCache) delete(key string) {
	if el, ok := c.items[key]; ok {
		delete(c.items, key) //map delete builtin
		c.order.Remove(el)
	}
}
