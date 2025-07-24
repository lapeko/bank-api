package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/v1/utils"
	apiUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

var service *userService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = &userService{store: store}
	g := router.Group(path)

	g.GET("/", listUsersHandler)
	g.GET("/:id", getUserByIdHandler)
	g.PATCH("/:id/email", updateUserEmailHandler)
	g.PATCH("/:id/fullname", updateUserFullNameHandler)
	g.PATCH("/:id/password", updateUserPasswordHandler)
	g.DELETE("/:id", deleteUserHandler)
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
	var uriId utils.UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	user, err := service.getUserById(ctx, uriId.ID)
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
	var uriId utils.UriId
	var body updateUserEmailRequest

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&body); err != nil {
		utils.SendError(ctx, err)
		return
	}

	user, err := service.updateUserEmail(ctx, uriId.ID, body.NewEmail)

	if err != nil {
		// TODO handle not unique email
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, user)
}

func updateUserFullNameHandler(ctx *gin.Context) {
	var uriId utils.UriId
	var body updateUserFullNameRequest

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&body); err != nil {
		utils.SendError(ctx, err)
		return
	}

	user, err := service.updateUserFullName(ctx, uriId.ID, body.NewFullName)

	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, user)
}

func updateUserPasswordHandler(ctx *gin.Context) {
	var uriId utils.UriId
	var body updateUserPasswordRequest

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&body); err != nil {
		utils.SendError(ctx, err)
		return
	}
	hash, err := apiUtils.HashPassword(body.NewPassword)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	user, err := service.updateUserPassword(ctx, uriId.ID, hash)

	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, user)
}

func deleteUserHandler(ctx *gin.Context) {
	var uriId utils.UriId

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}

	err := service.deleteUser(ctx, uriId.ID)

	if err != nil {
		// TODO recheck
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, uriId)
}
