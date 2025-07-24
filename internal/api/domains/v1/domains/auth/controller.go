package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/v1/utils"
	apiUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	internalUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

var service authService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = authService{store: store}
	g := router.Group(path)

	g.POST("/signup", signupHandler)
	g.POST("/signin", signinHandler)
	g.POST("/refresh", refreshHandler)
}

func signupHandler(ctx *gin.Context) {
	var usr createUserRequest
	if err := ctx.ShouldBind(&usr); err != nil {
		utils.SendError(ctx, err)
		return
	}
	hash, err := apiUtils.HashPassword(usr.Password)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	params := db.CreateUserParams{FullName: usr.FullName, Email: usr.Email, HashedPassword: hash}
	res, err := service.createUser(ctx, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == internalUtils.PgErrCodeUniqueViolation {
				utils.SendError(ctx, &authClientError{emailDuplicate})
				return
			}
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccessWithStatusCode(ctx, res, http.StatusCreated)
}

func signinHandler(ctx *gin.Context) {
	var req signinRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	tokens, err := service.signIn(ctx, req)
	var targetErr *authClientError
	if err != nil {
		if errors.As(err, &targetErr) {
			utils.SendError(ctx, err)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, tokens)
}

func refreshHandler(ctx *gin.Context) {
	var req refreshTokenRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	res, err := service.refreshToken(ctx, req)
	if err != nil {
		var target *authClientError
		if errors.As(err, &target) {
			utils.SendError(ctx, err)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, res)
}
