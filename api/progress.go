package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "progress.me-api/db/sql/sqlc"
)

type Progress struct {
	RangeValue    string `json:"range_value"`
	ProgressValue int64  `json:"progress_value"`
}

type CreateChartWithProgressesReq struct {
	UserId       uuid.UUID  `json:"user_id"`
	ProgressName string     `json:"progress_name"`
	RangeType    string     `json:"range_type"`
	ProgressData []Progress `json:"progress_data"`
}

type CreateChartWithProgressesRes struct {
	ChartID      uuid.UUID     `json:"chart_id"`
	UserID       uuid.UUID     `json:"user_id"`
	ProgressData []db.Progress `json:"progress_data"`
}

func (server *Server) CreateChartWithProgresses(ctx *gin.Context) {
	var req CreateChartWithProgressesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// insert into chart table
	chartData := db.CreateChartParams{
		UserID: uuid.NullUUID{
			UUID:  req.UserId,
			Valid: true,
		},
		RangeType:    db.Range(req.RangeType),
		ProgressName: req.ProgressName,
	}
	chart, err := server.store.CreateChart(ctx, chartData)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// insert into progress table using chart ID from above result
	for _, prog := range req.ProgressData {
		data := db.CreateProgressParams{
			ChartID: uuid.NullUUID{
				UUID:  chart.ID,
				Valid: true,
			},
			RangeValue:    prog.RangeValue,
			ProgressValue: prog.ProgressValue,
		}

		err := server.store.CreateProgress(ctx, data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// get all progresses from progress table using chart ID from above result
	progresses, err := server.store.GetProgressByChartID(ctx, uuid.NullUUID{UUID: chart.ID, Valid: true})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, CreateChartWithProgressesRes{
		ChartID:      chart.ID,
		UserID:       chart.UserID.UUID,
		ProgressData: progresses,
	})
}
