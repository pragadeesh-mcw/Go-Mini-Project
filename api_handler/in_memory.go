package api_handler

import (
	"net/http"
	"time"

	"github.com/pragadeesh-mcw/Go-Mini-Project/in_memory"

	"github.com/gin-gonic/gin"
)

// FOR IN_MEMORY ONLY INTERACTION
func SetupInMemoryRoutes(r *gin.Engine, cache in_memory.Cache) {
	r.GET("/inmemory/:key", func(c *gin.Context) {
		key := c.Param("key") //extract key from request
		value, found := cache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.GET("/inmemory", func(c *gin.Context) {
		items := cache.GetAll()
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.POST("/inmemory", func(c *gin.Context) {
		var json struct {
			Key        string `json:"key"`
			Value      string `json:"value"`
			Expiration int64  `json:"expiration"`
		}
		if err := c.ShouldBindJSON(&json); err != nil { //bind json with cacheItem struct
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cache.Set(json.Key, json.Value, time.Duration(json.Expiration)*time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully set value"})
	})

	r.DELETE("/inmemory/:key", func(c *gin.Context) {
		key := c.Param("key") //extract key from request
		if cache.Delete(key) {
			c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted key"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "Key is not found"})
		}
	})

	r.DELETE("/inmemory", func(c *gin.Context) {
		if cache.DeleteAll() {
			c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted all keys"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "No keys found"})
		}
	})

}
