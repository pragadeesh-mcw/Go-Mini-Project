package main

import (
	api "unified/api_handler"
	"unified/in_memory"
	"unified/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	rcache := redis.NewCache("localhost:6379", "", 0, 3)
	inMemoryCache := in_memory.NewLRUCache(3, 60)

	r := gin.Default()
	api.SetupRedisRoutes(r, rcache)

	api.SetupInMemoryRoutes(r, inMemoryCache)

	r.Run(":8080")
}
