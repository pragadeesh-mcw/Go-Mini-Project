package api_handler

import (
	"net/http"
	"time"

	"unified/multicache"

	"github.com/gin-gonic/gin"
)

func SetupUnifiedRoutes(multiCache *multicache.MultiCache) *gin.Engine {
	r := gin.Default()
	//SET
	r.POST("/cache", func(c *gin.Context) {
		var req struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
			TTL   int         `json:"ttl"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ttl := time.Duration(req.TTL) * time.Second
		if err := multiCache.Set(req.Key, req.Value, ttl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Key-Value pair set successfully"})
	})
	//GET
	r.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, err := multiCache.Get(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}

		if value == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "Key not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})
	//GETALL
	r.GET("/cache", func(c *gin.Context) {
		values, err := multiCache.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"values": values})
	})
	//DELETE
	r.DELETE("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		if err := multiCache.Delete(key); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Key deleted successfully"})
	})
	//DELETEALL
	r.DELETE("/cache", func(c *gin.Context) {
		if err := multiCache.DeleteAll(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "All keys deleted successfully"})
	})

	return r
}
