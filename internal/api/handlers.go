package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR PLN"`
}

func (a *Api) createAccount(ctx *gin.Context) {
	req := &createAccountRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	params := repository.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	acc, err := a.store.CreateAccount(context.Background(), params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusCreated, genOkBody(acc))
}

type getAccountsRequest struct {
	Page int32 `form:"page" json:"page" binding:"required,min=1"`
	Size int32 `form:"size" json:"size" binding:"required,min=5,max=50"`
}

func (a *Api) getAccounts(ctx *gin.Context) {
	req := &getAccountsRequest{}

	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	params := repository.ListAccountsParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	accounts, err := a.store.ListAccounts(context.Background(), params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusOK, genOkBody(accounts))
}

type getAccountByIdRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (a *Api) getAccount(ctx *gin.Context) {
	req := &getAccountByIdRequest{}

	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	accounts, err := a.store.GetAccount(context.Background(), req.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, genFailBody("user not found"))
			return
		}
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusOK, genOkBody(accounts))
}
