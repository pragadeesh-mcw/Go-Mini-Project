package api_handler

import (
	"net/http"
	"time"

	"unified/redis"

	"github.com/gin-gonic/gin"
)

func SetupRedisRoutes(r *gin.Engine, cache redis.Cache) {
	r.GET("/redis/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, found := cache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.GET("/redis/cache", func(c *gin.Context) {
		keys := cache.GetAll()
		c.JSON(http.StatusOK, gin.H{"keys": keys})
	})

	r.POST("/redis/cache", func(c *gin.Context) {
		var json struct {
			Key        string `json:"key"`
			Value      string `json:"value"`
			Expiration int64  `json:"expiration"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cache.Set(json.Key, json.Value, time.Duration(json.Expiration)*time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully set value"})
	})

	r.DELETE("/redis/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		cache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted key"})
	})

	r.DELETE("/redis/cache", func(c *gin.Context) {
		cache.DeleteAll()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted all keys"})
	})
}
