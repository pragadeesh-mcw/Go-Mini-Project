package main

import (
	"flag"
	"log"
	api "unified/api_handler"
	"unified/in_memory"
	"unified/multicache"
	"unified/redis_cache"

	"github.com/gin-gonic/gin"
)

func main() {
	//set max capacity
	var maxCacheCapacity int
	flag.IntVar(&maxCacheCapacity, "cache-capacity", 3, "Maximum capacity of the cache")
	flag.Parse()
	//initiate redis and in-memory
	inMemoryCache := in_memory.NewLRUCache(maxCacheCapacity, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, maxCacheCapacity)
	multiCache := multicache.NewMultiCache(inMemoryCache, redisCache)
	//setup unified api
	r1 := gin.Default()
	api.SetupInMemoryRoutes(r1, inMemoryCache)
	api.SetupRedisRoutes(r1, redisCache)
	r := api.SetupUnifiedRoutes(multiCache)

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

	// Block forever
	select {}
}
