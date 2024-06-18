package main

import (
	api "unified/api_handler"
	"unified/in_memory"
	"unified/multicache"
	"unified/redis"
)

func main() {
	inMemoryCache := in_memory.NewLRUCache(100, 60)
	redisCache := redis.NewCache("localhost:6379", "", 0, 100)
	multiCache := multicache.NewMultiCache(inMemoryCache, redisCache)

	r := api.SetupRouter(multiCache)
	r.Run(":8080")
}
