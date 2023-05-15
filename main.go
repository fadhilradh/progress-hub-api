package main

import (
	"log"

	"github.com/redis/go-redis/v9"
	"progress.me-api/api"
	db "progress.me-api/db/redis"
)

var redisOpt = &redis.Options{
	Addr:     "localhost:6379",
	Password: "authentication",
	DB:       0,
}

func main() {
	rdb := db.NewRedis(redisOpt)
	server := api.NewServer(rdb)

	err := server.Start("0.0.0.0:8080")

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
