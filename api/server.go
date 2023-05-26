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
	// router.Group("/api/v1")
	// router.BasePath()
	router.Use(cors.Default())
	router.GET("/api/v1/health", testFunc)

	router.GET("/api/v1/charts/:id", server.GetChartByID)

	router.PATCH("/api/v1/progresses", server.BulkUpdateProgress)
	router.POST("/api/v1/progresses", server.CreateProgress)
	router.DELETE("/api/v1/progresses/:id", server.DeleteProgressByID)

	router.POST("/api/v1/chart-progresses", server.CreateChartWithProgresses)
	router.GET("/api/v1/chart-progresses/:user_id", server.ListChartProgressByUserId)

	router.POST("/api/v1/cache/progress", server.CacheProgress)
	router.POST("/api/v1/users", server.createUser)

	router.POST("/api/v1/login", server.LoginUser)

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
