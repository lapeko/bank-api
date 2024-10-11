package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func setUpUsers(r *gin.Engine) {
	users := r.Group("/users")

	users.POST("/", createUser)
}

type createUserRequest struct {
	FullName string `json:"fullName" binding:"alpha"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"min=6"`
}

type createUserResponse struct {
	ID                int64     `json:"id"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangesAt time.Time `json:"password_changes_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func createUser(ctx *gin.Context) {
	a := GetApi()
	req := createUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, genFailBody(err))
		return
	}

	user, err := a.store.CreateUser(context.Background(), repository.CreateUserParams{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	})

	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			switch e.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, genFailBody(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, genFailBody(err))
		return
	}

	userWithoutPassword := &createUserResponse{
		ID:                user.ID,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangesAt: user.PasswordChangesAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, genOkBody(userWithoutPassword))
}
