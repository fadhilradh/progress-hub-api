package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "progress.me-api/db/sql/sqlc"
)

type CreateProgressReq struct {
	ChartID       uuid.UUID `json:"chart_id"`
	RangeValue    string    `json:"range_value"`
	ProgressValue int64     `json:"progress_value"`
}

type CacheProgressReq struct {
	ID   string        `json:"id"`
	Data []db.Progress `json:"data"`
}

func (i CacheProgressReq) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (server *Server) CreateProgress(ctx *gin.Context) {
	var req CreateProgressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data := db.CreateProgressParams{
		ChartID: uuid.NullUUID{
			UUID:  req.ChartID,
			Valid: true,
		},
		RangeValue:    req.RangeValue,
		ProgressValue: req.ProgressValue,
	}

	progress, err := server.store.CreateProgress(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, progress)
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
