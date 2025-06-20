package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "progress.me-api/db/sql/sqlc"
)

type ChartType string

const (
	ChartBar  ChartType = "bar"
	ChartArea ChartType = "area"
	ChartLine ChartType = "line"
)

type BarChartType string

const (
	Horizontal BarChartType = "horizontal"
	Vertical   BarChartType = "vertical"
)

type Progress struct {
	RangeValue    *string `json:"range_value"`
	ProgressValue *int64  `json:"progress_value"`
	ProgressNo    *int32  `json:"progress_no"`
}

type CreateChartWithProgressesReq struct {
	UserId       *string    `json:"user_id", binding:"required"`
	ProgressName *string      `json:"progress_name", binding:"required"`
	RangeType    *string      `json:"range_type", binding:"required"`
	ProgressData []Progress   `json:"progress_data", binding:"required"`
	ChartColor   *string      `json:"chart_color", binding:"required"`
	ChartType    ChartType    `json:"chart_type", binding:"required"`
	BarChartType BarChartType `json:"bar_chart_type"`
}

type CreateChartWithProgressesRes struct {
	ChartID      uuid.UUID     `json:"chart_id"`
	UserID       string     `json:"user_id"`
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
		UserID: req.UserId,
		RangeType:    req.RangeType,
		ProgressName: req.ProgressName,
		Colors:       req.ChartColor,
		ChartType:    (*string)(&req.ChartType),
		BarChartType: (*string)(&req.BarChartType),
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
			ProgressNo:    prog.ProgressNo,
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
		UserID:       *chart.UserID,
		ProgressData: progresses,
	})
}

type ProgressData struct {
	ProgressID    uuid.UUID `json:"progress_id"`
	RangeValue    *string   `json:"range_value"`
	ProgressValue *int64    `json:"progress_value"`
	ProgressNo    *int32    `json:"progress_no"`
}

type GetChartsByUserIdRes struct {
	ChartID      uuid.UUID      `json:"chart_id"`
	ChartColor   *string        `json:"chart_color"`
	ChartType    ChartType      `json:"chart_type"`
	BarChartType BarChartType   `json:"bar_chart_type"`
	RangeType    *string        `json:"range_type"`
	ProgressName *string        `json:"progress_name"`
	ProgressData []ProgressData `json:"progress_data"`
}

type GetAccountByIDReq struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

func (server *Server) ListChartProgressByUserId(ctx *gin.Context) {
	userID := (ctx.Param("user_id"))
	charts, err := server.store.ListChartProgressByUserId(ctx, &userID)
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
				ChartType:    ChartType(*ch.ChartType),
				BarChartType: BarChartType(*ch.BarChartType),
			})
			currChartID = ch.ChartID
			if i != 0 {
				currChartIdx++
			}
			chartRes[currChartIdx].ProgressData = append(chartRes[currChartIdx].ProgressData, ProgressData{
				ProgressID:    ch.ProgressID,
				RangeValue:    ch.RangeValue,
				ProgressValue: ch.ProgressValue,
				ProgressNo:    ch.ProgressNo,
			})

		} else {
			chartRes[currChartIdx].ProgressData = append(chartRes[currChartIdx].ProgressData, ProgressData{
				ProgressID:    ch.ProgressID,
				RangeValue:    ch.RangeValue,
				ProgressValue: ch.ProgressValue,
				ProgressNo:    ch.ProgressNo,
			})
		}

	}

	ctx.JSON(http.StatusOK, chartRes)
}

type GetChartByIDRes struct {
	ChartID      uuid.UUID      `json:"chart_id"`
	ChartColor   *string        `json:"chart_color"`
	ChartType    ChartType      `json:"chart_type"`
	BarChartType BarChartType   `json:"bar_chart_type"`
	RangeType    *string        `json:"range_type"`
	ProgressName *string        `json:"progress_name"`
	ProgressData []ProgressData `json:"progress_data"`
}

func (server *Server) GetChartByID(ctx *gin.Context) {
	chartID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	chart, err := server.store.GetChartByID(ctx, chartID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	progress, err := server.store.GetProgressByChartID(ctx, uuid.NullUUID{
		UUID:  chart.ID,
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

	response := GetChartByIDRes{
		ChartID:      chart.ID,
		ChartColor:   chart.Colors,
		ChartType:    ChartType(*chart.ChartType),
		BarChartType: BarChartType(*chart.BarChartType),
		RangeType:    chart.RangeType,
		ProgressName: chart.ProgressName,
	}

	for _, prog := range progress {
		response.ProgressData = append(response.ProgressData, ProgressData{
			ProgressID:    prog.ID,
			RangeValue:    prog.RangeValue,
			ProgressValue: prog.ProgressValue,
			ProgressNo:    prog.ProgressNo,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

type UpdateChartReq struct {
	ProgressName *string `json:"progress_name"`
	ChartColor   *string `json:"chart_color"`
	ChartType    *string `json:"chart_type"`
	BarChartType *string `json:"bar_chart_type"`
	RangeType    *string `json:"range_type"`
}

func (server *Server) UpdateChartByID(ctx *gin.Context) {
	chartID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateChartReq
	log.Print("req", req)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := db.UpdateChartParams{
		ID:           chartID,
		RangeType:    req.RangeType,
		ProgressName: req.ProgressName,
		Colors:       req.ChartColor,
		ChartType:    req.ChartType,
		BarChartType: req.BarChartType,
	}

	log.Print("param", param)

	if err := server.store.UpdateChart(ctx, param); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Chart updated successfully",
	})
}
