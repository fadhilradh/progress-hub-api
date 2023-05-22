package api

import (
	"database/sql"
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

type ProgressData struct {
	ProgressID    uuid.UUID `json:"progress_id"`
	RangeValue    string    `json:"range_value"`
	ProgressValue int64     `json:"progress_value"`
}

type GetChartsByUserIdRes struct {
	ChartID      uuid.UUID      `json:"chart_id"`
	RangeType    db.Range       `json:"range_type"`
	ProgressName string         `json:"progress_name"`
	ProgressData []ProgressData `json:"progress_data"`
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

type GetAccountByIDReq struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetChartProgressByUserId(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	charts, err := server.store.GetChartProgressByUserId(ctx, uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var chartRes []GetChartsByUserIdRes
	currChartID := charts[0].ChartID
	currChartIdx := 0
	// Iterate over charts : Make one object for each chartID
	for i, ch := range charts {
		// if the first index or
		// if the chartID is different from the previous one, create new chart object
		if i == 0 || ch.ChartID != currChartID {
			chartRes = append(chartRes, GetChartsByUserIdRes{
				ChartID:      ch.ChartID,
				RangeType:    ch.RangeType,
				ProgressName: ch.ProgressName,
			})
			currChartID = ch.ChartID
			if i != 0 {
				currChartIdx++
			}
			chartRes[currChartIdx].ProgressData = append(chartRes[currChartIdx].ProgressData, ProgressData{
				ProgressID:    ch.ProgressID,
				RangeValue:    ch.RangeValue,
				ProgressValue: ch.ProgressValue,
			})

		} else {
			chartRes[currChartIdx].ProgressData = append(chartRes[currChartIdx].ProgressData, ProgressData{
				ProgressID:    ch.ProgressID,
				RangeValue:    ch.RangeValue,
				ProgressValue: ch.ProgressValue,
			})
		}

	}

	ctx.JSON(http.StatusOK, chartRes)
}
