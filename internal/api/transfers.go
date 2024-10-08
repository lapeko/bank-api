package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	"net/http"
)

func setUpTransfers(r *gin.Engine) {
	accounts := r.Group("/transfers")

	accounts.POST("/", createTransfer)
}

type createTransferRequest struct {
	AccountFrom int64 `json:"accountFrom" binding:"required,min=0"`
	AccountTo   int64 `json:"accountTo" binding:"required,min=0"`
	Amount      int64 `json:"amount" binding:"required,gt=0"`
}

func createTransfer(ctx *gin.Context) {
	a := GetApi()
	req := createTransferRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	accFrom, err := a.store.GetAccount(context.Background(), req.AccountFrom)

	if err != nil {
		ctx.JSON(http.StatusNotFound, genFailBody("accountFrom not found"))
		return
	}

	accTo, err := a.store.GetAccount(context.Background(), req.AccountTo)

	if err != nil {
		ctx.JSON(http.StatusNotFound, genFailBody("accountTo not found"))
		return
	}

	if accFrom.Currency != accTo.Currency {
		ctx.JSON(http.StatusBadRequest, genFailBody("accounts have different currencies"))
		return
	}

	if accFrom.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, genFailBody("insufficient funds"))
		return
	}

	res, err := a.store.TransferTX(context.Background(), repository.CreateTransferParams(req))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusCreated, genOkBody(res))
}
