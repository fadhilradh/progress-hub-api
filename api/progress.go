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
