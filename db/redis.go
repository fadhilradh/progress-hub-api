package db

import (
	"github.com/go-redis/redis"
)

func NewRedis(opt *redis.Options) *redis.Client {
	rdb := redis.NewClient(opt)

	return rdb

}
