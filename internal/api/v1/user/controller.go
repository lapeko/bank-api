package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	rootUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

var service *userService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = &userService{store: store}
	g := router.Group(path)

	g.POST("/", createUserHandler)
	g.GET("/", listUsersHandler)
	g.GET("/:id", getUserByIdHandler)
	g.PATCH("/:id/email", updateUserEmailHandler)
	g.PATCH("/:id/fullname", updateUserFullNameHandler)
	g.PATCH("/:id/password", updateUserPasswordHandler)
	g.DELETE("/:id", deleteUserHandler)
}

func createUserHandler(ctx *gin.Context) {
	var usr createUserRequest
	if err := ctx.ShouldBind(&usr); err != nil {
		utils.SendError(ctx, err)
		return
	}
	hash, err := rootUtils.HashPassword(usr.Password)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	params := db.CreateUserParams{FullName: usr.FullName, Email: usr.Email, HashedPassword: hash}
	usrRes, err := service.createUser(ctx, params)
	if err != nil {
		// TODO handle email taken
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccessWithStatusCode(ctx, usrRes, http.StatusCreated)
}

func listUsersHandler(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	res, err := service.listUsers(ctx, req)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, res)
}

func getUserByIdHandler(ctx *gin.Context) {
	var param userUriIdRequest
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.SendError(ctx, err)
		return
	}
	user, err := service.getUserById(ctx, param.ID)
	if err != nil {
		// TODO recheck
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, user)
}

func updateUserEmailHandler(ctx *gin.Context) {
	var uriPath userUriIdRequest
	var body updateUserEmailRequest

	if err := ctx.ShouldBindUri(&uriPath); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&body); err != nil {
		utils.SendError(ctx, err)
		return
	}

	user, err := service.updateUserEmail(ctx, uriPath.ID, body.NewEmail)

	if err != nil {
		// TODO handle not unique email
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, user)
}

func updateUserFullNameHandler(ctx *gin.Context) {
	var uriPath userUriIdRequest
	var body updateUserFullNameRequest

	if err := ctx.ShouldBindUri(&uriPath); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&body); err != nil {
		utils.SendError(ctx, err)
		return
	}

	user, err := service.updateUserFullName(ctx, uriPath.ID, body.NewFullName)

	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, user)
}

func updateUserPasswordHandler(ctx *gin.Context) {
	var uriPath userUriIdRequest
	var body updateUserPasswordRequest

	if err := ctx.ShouldBindUri(&uriPath); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&body); err != nil {
		utils.SendError(ctx, err)
		return
	}
	hash, err := rootUtils.HashPassword(body.NewPassword)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	user, err := service.updateUserPassword(ctx, uriPath.ID, hash)

	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, user)
}

func deleteUserHandler(ctx *gin.Context) {
	var uriPath userUriIdRequest

	if err := ctx.ShouldBindUri(&uriPath); err != nil {
		utils.SendError(ctx, err)
		return
	}

	err := service.deleteUser(ctx, uriPath.ID)

	if err != nil {
		// TODO recheck
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, uriPath)
}
