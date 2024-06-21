package main

import (
	"log"
	api "unified/api_handler"
	"unified/in_memory"
	"unified/multicache"
	"unified/redis_cache"

	"github.com/gin-gonic/gin"
)

func main() {
	//initiate redis and in-memory
	inMemoryCache := in_memory.NewLRUCache(3, 60)
	redisCache := redis_cache.NewCache("localhost:6379", "", 0, 3)
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
