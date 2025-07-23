package transfer

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

var service transferService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = transferService{store: store}
	g := router.Group(path)

	g.POST("/", transferHandler)
	g.POST("/in", transferInHandler)
	g.POST("/out", transferOutHandler)
	g.GET("/", listTransfersHandler)
	g.GET("/account/:id", listTransfersByAccountHandler)
}

func transferHandler(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := service.transferMoney(ctx, req); err != nil {
		var transferClientError *db.TransferClientError
		if errors.As(err, &transferClientError) {
			utils.SendError(ctx, err)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, struct{}{})
}

func transferInHandler(ctx *gin.Context) {
	var req externalTransferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := service.transferIn(ctx, req); err != nil {
		var target *db.TransferClientError
		if errors.As(err, &target) {
			utils.SendError(ctx, err)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, struct{}{})
}

func transferOutHandler(ctx *gin.Context) {
	var req externalTransferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	acc, err := service.transferOut(ctx, req)
	if err != nil {
		var target *db.TransferClientError
		if errors.As(err, &target) {
			utils.SendError(ctx, err)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, acc)
}

func listTransfersHandler(ctx *gin.Context) {
	var paginationParams listTransfersRequest

	if err := ctx.ShouldBind(&paginationParams); err != nil {
		utils.SendError(ctx, err)
		return
	}

	response, err := service.listTransfers(ctx, paginationParams)

	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, response)
}

func listTransfersByAccountHandler(ctx *gin.Context) {
	var uriId utils.UriId
	var paginationParams listTransfersRequest

	if err := ctx.ShouldBind(&paginationParams); err != nil {
		utils.SendError(ctx, err)
		return
	}

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}

	response, err := service.listTransfersByAccount(ctx, paginationParams, uriId.ID)

	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(ctx, response)
}
