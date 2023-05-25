package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"progress.me-api/api"
	db "progress.me-api/db/sql/sqlc"
	"progress.me-api/util"
)

var redisOpt = &redis.Options{
	Addr:     "localhost:6379",
	Password: "authentication",
	DB:       0,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DSN)
	if err != nil {
		log.Fatal("Can't connect to db:", err)
		return
	}
	store := db.NewStore(conn)
	// rdb := reddb.NewRedis(redisOpt)
	server := api.NewServer(nil, store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
