package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "progress.me-api/db/sql/sqlc"
)

type CreateChartReq struct {
	UserID       uuid.UUID `json:"user_id"`
	RangeType    string    `json:"range_type"`
	ProgressName string    `json:"progress_name"`
	ChartColor   string    `json:"chart_color"`
}

type ProgressData struct {
	ProgressID    uuid.UUID `json:"progress_id"`
	RangeValue    string    `json:"range_value"`
	ProgressValue int64     `json:"progress_value"`
}

type GetChartsByUserIdRes struct {
	ChartID      uuid.UUID      `json:"chart_id"`
	ChartColor   string         `json:"chart_color"`
	RangeType    db.Range       `json:"range_type"`
	ProgressName string         `json:"progress_name"`
	ProgressData []ProgressData `json:"progress_data"`
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
				ChartColor:   ch.ChartColor,
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

type Progress struct {
	RangeValue    string `json:"range_value"`
	ProgressValue int64  `json:"progress_value"`
}

type CreateChartWithProgressesReq struct {
	UserId       uuid.UUID  `json:"user_id"`
	ProgressName string     `json:"progress_name"`
	RangeType    string     `json:"range_type"`
	ProgressData []Progress `json:"progress_data"`
	ChartColor   string     `json:"chart_color"`
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
		Colors:       req.ChartColor,
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
