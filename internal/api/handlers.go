package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR PLN"`
}

func (a *Api) createAccount(ctx *gin.Context) {
	accReq := &createAccountRequest{}
	if err := ctx.ShouldBindJSON(accReq); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailRes(err))
		return
	}

	params := repository.CreateAccountParams{
		Owner:    accReq.Owner,
		Currency: accReq.Currency,
		Balance:  0,
	}
	acc, err := a.store.CreateAccount(context.Background(), params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, genFailRes(err))
		return
	}

	ctx.JSON(http.StatusCreated, genOkRes(acc))
}
