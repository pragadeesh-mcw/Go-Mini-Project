package cache

import (
	api "github.com/pragadeesh-mcw/Go-Mini-Project/api_handler"

	"github.com/pragadeesh-mcw/Go-Mini-Project/in_memory"
	"github.com/pragadeesh-mcw/Go-Mini-Project/multicache"
	"github.com/pragadeesh-mcw/Go-Mini-Project/redis_cache"

	"github.com/gin-gonic/gin"
)

func entry() *gin.Engine {
	//initiate redis and in-memory
	inMemoryCache := in_memory.NewLRUCache(3, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, 3)
	multiCache := multicache.NewMultiCache(inMemoryCache, redisCache)
	//setup unified api

	r := api.SetupRouter(multiCache)
	api.SetupRedisRoutes(r, redisCache)

	api.SetupInMemoryRoutes(r, inMemoryCache)
	return r
}
