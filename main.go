package main

import (
	api "unified/api_handler"
	"unified/in_memory"
	"unified/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	redis.InitRedis("localhost:6379", "", 0)
	addr := "localhost:6379"
	password := ""
	db := 0

	redisCache := redis.NewCache(addr, password, db)
	inMemoryCache := in_memory.NewLRUCache(3, 60)

	r := gin.Default()

	api.SetupRedisRoutes(r, redisCache)
	api.SetupInMemoryRoutes(r, inMemoryCache)

	r.Run(":8080")
}
