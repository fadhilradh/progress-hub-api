package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	router.Use(cors.Default())
	router.GET("/", testFunc)
	router.POST("/progress", server.createProgress)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run()
}
