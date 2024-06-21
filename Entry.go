package cache

import (
	"log"

	"github.com/gin-gonic/gin"
	api "github.com/pragadeesh-mcw/Go-Mini-Project/api_handler"
	"github.com/pragadeesh-mcw/Go-Mini-Project/in_memory"
	"github.com/pragadeesh-mcw/Go-Mini-Project/multicache"
	"github.com/pragadeesh-mcw/Go-Mini-Project/redis_cache"
)

// GinEngines holds the Gin engines r and r1
type GinEngines struct {
	R  *gin.Engine
	R1 *gin.Engine
}

// Entry initializes and returns the Gin engines r and r1
func Entry() *GinEngines {
	// Initialize caches
	inMemoryCache := in_memory.NewLRUCache(3, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, 3)
	multiCache := multicache.NewMultiCache(inMemoryCache, redisCache)

	// Setup unified API
	r1 := gin.Default()
	api.SetupInMemoryRoutes(r1, inMemoryCache)
	api.SetupRedisRoutes(r1, redisCache)

	r := gin.Default()
	api.SetupUnifiedRoutes(multiCache)

	// Run servers concurrently
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server on :8080: %v", err)
		}
	}()

	go func() {
		if err := r1.Run(":8081"); err != nil {
			log.Fatalf("Failed to run server on :8081: %v", err)
		}
	}()

	return &GinEngines{
		R:  r,
		R1: r1,
	}
}
