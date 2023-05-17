package db

import (
	"github.com/redis/go-redis/v9"
)

func NewRedis(opt *redis.Options) *redis.Client {
	rdb := redis.NewClient(opt)

	return rdb

}
