package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

var service *accountService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = &accountService{store: store}
	g := router.Group(path)

	g.POST("/", createAccountHandler)
	g.GET("/", listAccountsHandler)
	g.GET("/:id", getAccountByIdHandler)
	g.DELETE("/:id", deleteAccountHandler)
}

func createAccountHandler(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	account, err := service.createAccount(ctx, &req)
	if err != nil {
		// TODO handle duplicated key errors
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccessWithStatusCode(ctx, account, http.StatusCreated)
}

func listAccountsHandler(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	res, err := service.listAccounts(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, res)
}

func getAccountByIdHandler(ctx *gin.Context) {
	var uriId utils.UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	acc, err := service.getAccountById(ctx, uriId.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, acc)
}

func deleteAccountHandler(ctx *gin.Context) {
	var uriId utils.UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := service.deleteAccountById(ctx, uriId.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, uriId)
}
