package api_handler

import (
	"net/http"
	"time"
	"unified/redis"

	"github.com/gin-gonic/gin"
)

var cacheInstance *redis.Cache

func SetupRedisRoutes(router *gin.Engine, cache *redis.Cache) {
	cacheInstance = cache

	router.POST("/redis", setHandler)
	router.GET("/redis/get", getHandler)
	router.GET("/redis", getAllHandler)
	router.DELETE("/redis/:key", deleteHandler)
	router.DELETE("/redis", deleteAllHandler)
}

type SetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl" binding:"required"`
}

func setHandler(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cacheInstance.Set(req.Key, req.Value, time.Duration(req.TTL)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Key set successfully"})
}

func getHandler(c *gin.Context) {
	key := c.Query("key")
	value, err := cacheInstance.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func getAllHandler(c *gin.Context) {
	values, err := cacheInstance.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, values)
}

func deleteHandler(c *gin.Context) {
	key := c.Query("key")
	err := cacheInstance.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
}

func deleteAllHandler(c *gin.Context) {
	err := cacheInstance.DeleteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All keys deleted successfully"})
}
