package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "progress.me-api/db/sql/sqlc"
)

type CacheProgressReq struct {
	ID   string        `json:"id"`
	Data []db.Progress `json:"data"`
}

func (i CacheProgressReq) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (server *Server) CacheProgress(ctx *gin.Context) {
	var req CacheProgressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var c = context.Background()

	err := server.redis.Set(c, req.ID, req, 0).Err()
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	val, err := server.redis.Get(c, req.ID).Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Print(err)
		return
	}
	var raw CacheProgressReq
	if err := json.Unmarshal([]byte(val), &raw); err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, raw)
}
