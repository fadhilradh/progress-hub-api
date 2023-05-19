package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "progress.me-api/db/sql/sqlc"
)

type CreateChartReq struct {
	UserID       uuid.UUID `json:"user_id"`
	RangeType    string    `json:"range_type"`
	ProgressName string    `json:"progress_name"`
}

func (server *Server) CreateChart(ctx *gin.Context) db.Chart {
	var req CreateChartReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return db.Chart{}
	}

	data := db.CreateChartParams{
		UserID: uuid.NullUUID{
			UUID:  req.UserID,
			Valid: true,
		},
		RangeType:    db.Range(req.RangeType),
		ProgressName: req.ProgressName,
	}

	progress, err := server.store.CreateChart(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Chart{}
	}

	return progress
}
