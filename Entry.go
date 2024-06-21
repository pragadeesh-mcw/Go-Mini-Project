package cache

import (
	api "github.com/pragadeesh-mcw/Go-Mini-Project/api_handler"

	"github.com/pragadeesh-mcw/Go-Mini-Project/in_memory"
	"github.com/pragadeesh-mcw/Go-Mini-Project/multicache"
	"github.com/pragadeesh-mcw/Go-Mini-Project/redis_cache"

	"github.com/gin-gonic/gin"
)

type GinEngines struct {
	R  *gin.Engine
	R1 *gin.Engine
}

func Entry() *GinEngines {
	inMemoryCache := in_memory.NewLRUCache(3, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, 3)
	multiCache := multicache.NewMultiCache(inMemoryCache, redisCache)

	// Setup unified API
	r1 := gin.Default()
	api.SetupInMemoryRoutes(r1, inMemoryCache)
	api.SetupRedisRoutes(r1, redisCache)

	r := gin.Default()
	api.SetupUnifiedRoutes(multiCache)

	return &GinEngines{
		R:  r,
		R1: r1,
	}
}
