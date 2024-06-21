package cache

import (
	api "github.com/pragadeesh-mcw/Go-Mini-Project/api_handler"

	"github.com/pragadeesh-mcw/Go-Mini-Project/in_memory"
	"github.com/pragadeesh-mcw/Go-Mini-Project/multicache"
	"github.com/pragadeesh-mcw/Go-Mini-Project/redis_cache"

	"github.com/gin-gonic/gin"
)

func Entry() *gin.Engine {
	inMemoryCache := in_memory.NewLRUCache(3, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, 3)
	multiCache := multicache.NewMultiCache(inMemoryCache, redisCache)

	// Setup unified API
	r:=api.SetupUnifiedRoutes(multiCache)
	return r
}
