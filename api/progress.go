package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateProgressReq struct {
	ID           string `json:"id"`
	DurationType string `json:"duration_type"` // month, year
	Value        int    `json:"value"`
	Date         string `json:"date,omitempty"` // used if duration type is date
}

func (i CreateProgressReq) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (server *Server) createProgress(ctx *gin.Context) {
	var req CreateProgressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var c = context.Background()

	err := server.redis.Set(c, "progress", req, 0).Err()
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	val, err := server.redis.Get(c, "progress").Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Print(err)
		return
	}
	var raw CreateProgressReq
	if err := json.Unmarshal([]byte(val), &raw); err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, raw)
}
