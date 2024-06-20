package main

import (
	api "unified/api_handler"
	"unified/in_memory"
	"unified/multicache"
	"unified/redis_cache"
)

func Entry() *gin.Engine {
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
