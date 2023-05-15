package main

import (
	"log"

	"github.com/go-redis/redis"
	"progress.me-api/api"
	"progress.me-api/db"
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
