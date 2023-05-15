package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Server struct {
	redis  *redis.Client
	router *gin.Engine
}

func testFunc(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func NewServer(redis *redis.Client) *Server {
	server := &Server{
		redis: redis,
	}

	router := gin.Default()

	router.GET("/", testFunc)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run()
}
