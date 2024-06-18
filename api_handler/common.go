package api_handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Cache interface {
	Get(key string) (string, bool)
	GetAll() interface{}
	Set(key, value string, expiration time.Duration)
	Delete(key string)
	DeleteAll()
}

func SetupUnifiedRoutes(r *gin.Engine, redisCache Cache, inMemoryCache Cache) {
	r.GET("/cache/redis/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, found := redisCache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found in Redis cache"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.GET("/cache/redis", func(c *gin.Context) {
		items := redisCache.GetAll()
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.POST("/cache/redis", func(c *gin.Context) {
		var json struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			TTL   int64  `json:"ttl"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		redisCache.Set(json.Key, json.Value, time.Duration(json.TTL)*time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully set value in Redis cache"})
	})

	r.DELETE("/cache/redis/:key", func(c *gin.Context) {
		key := c.Param("key")
		redisCache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted key from Redis cache"})
	})

	r.DELETE("/cache/redis", func(c *gin.Context) {
		redisCache.DeleteAll()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted all keys from Redis cache"})
	})

	r.GET("/cache/inmemory/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, found := inMemoryCache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found in in-memory cache"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.GET("/cache/inmemory", func(c *gin.Context) {
		items := inMemoryCache.GetAll()
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.POST("/cache/inmemory", func(c *gin.Context) {
		var json struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			TTL   int64  `json:"ttl"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		inMemoryCache.Set(json.Key, json.Value, time.Duration(json.TTL)*time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully set value in in-memory cache"})
	})

	r.DELETE("/cache/inmemory/:key", func(c *gin.Context) {
		key := c.Param("key")
		inMemoryCache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted key from in-memory cache"})
	})

	r.DELETE("/cache/inmemory", func(c *gin.Context) {
		inMemoryCache.DeleteAll()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted all keys from in-memory cache"})
	})
}
