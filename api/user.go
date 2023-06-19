package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "progress.me-api/db/sql/sqlc"
)

type Role string

const (
	RoleUser       Role = "user"
	RoleAdmin      Role = "admin"
	RoleSuperadmin Role = "superadmin"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data := db.CreateUserParams{
		Username: sql.NullString{
			String: req.Username,
			Valid:  true,
		},
		Password: sql.NullString{
			String: req.Password,
			Valid:  true,
		}, // encrypt this !
		Email: sql.NullString{
			String: req.Email,
			Valid:  true,
		},
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Role: sql.NullString{
			String: req.Role,
			Valid:  true,
		},
	}

	user, err := server.store.CreateUser(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) addUserFromClerk(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := server.store.AddClerkUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
