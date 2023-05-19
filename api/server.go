package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	db "progress.me-api/db/sql/sqlc"
)

type Server struct {
	redis  *redis.Client
	router *gin.Engine
	store  *db.Store
}

func testFunc(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func NewServer(redis *redis.Client, store *db.Store) *Server {
	server := &Server{
		redis: redis,
		store: store,
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", testFunc)
	router.POST("/chart", server.CreateChartWithProgresses)
	router.POST("/cache/progress", server.CacheProgress)
	router.POST("/users", server.createUser)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
