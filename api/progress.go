package api

import "github.com/gin-gonic/gin"

func (server *Server) createProgress(ctx *gin.Context) {
	server.redis.Set("test", "12345", 0)
}
