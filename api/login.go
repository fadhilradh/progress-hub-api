package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginUserReq struct {
	Username string `json:"username", binding:"required"`
	Password string `json:"password", binding:"required"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, sql.NullString{
		String: req.Username,
		Valid:  true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Password.String != req.Password {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
