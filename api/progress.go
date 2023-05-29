package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "progress.me-api/db/sql/sqlc"
)

type ProgressUpdate struct {
	ProgressID    uuid.UUID `json:"progress_id" binding:"required"`
	ProgressValue *int64    `json:"progress_value"`
	RangeValue    *string   `json:"range_value"`
	ProgressNo    *int32    `json:"progress_no"`
}

type BulkProgressReq struct {
	Progresses []ProgressUpdate `json:"progresses" binding:"required"`
}

func (server *Server) BulkUpdateProgress(ctx *gin.Context) {
	var req BulkProgressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Print(err)
		return
	}

	for _, prog := range req.Progresses {
		param := db.EditProgressByIDParams{
			ID:            prog.ProgressID,
			RangeValue:    prog.RangeValue,
			ProgressValue: prog.ProgressValue,
			ProgressNo:    prog.ProgressNo,
		}
		if err := server.store.EditProgressByID(ctx, param); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Progresses updated successfully",
	})
}

func (server *Server) DeleteProgressByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteProgressByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Progress deleted successfully",
	})
}

type CreateProgressReq struct {
	ChartID       uuid.UUID `json:"chart_id" binding:"required"`
	ProgressValue *int64    `json:"progress_value" binding:"required"`
	RangeValue    *string   `json:"range_value" binding:"required"`
	ProgressNo    *int32    `json:"progress_no" binding:"required"`
}

func (server *Server) CreateProgress(ctx *gin.Context) {
	var req CreateProgressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := db.CreateProgressParams{
		ChartID: uuid.NullUUID{
			UUID:  req.ChartID,
			Valid: true,
		},
		ProgressValue: req.ProgressValue,
		RangeValue:    req.RangeValue,
		ProgressNo:    req.ProgressNo,
	}

	err := server.store.CreateProgress(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Progress created successfully",
	})
}
