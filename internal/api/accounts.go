package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	"net/http"
)

func setUpAccounts(r *gin.Engine) {
	accounts := r.Group("/accounts")

	accounts.POST("/", createAccount)
	accounts.GET("/", getAccounts)
	accounts.GET("/:id", getAccount)
	accounts.PUT("/", updateAccount)
	accounts.DELETE("/:id", deleteAccount)
}

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func createAccount(ctx *gin.Context) {
	a := GetApi()
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

func getAccounts(ctx *gin.Context) {
	a := GetApi()
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

func getAccount(ctx *gin.Context) {
	a := GetApi()
	req := &getAccountByIdRequest{}

	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	account, err := a.store.GetAccount(context.Background(), req.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, genFailBody("user not found"))
			return
		}
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusOK, genOkBody(account))
}

type updateAccountRequest struct {
	ID      int64  `json:"id" binding:"required,min=1"`
	Balance *int64 `json:"balance"`
}

func updateAccount(ctx *gin.Context) {
	a := GetApi()
	req := updateAccountRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	if req.Balance == nil {
		ctx.JSON(http.StatusBadRequest, genFailBody("balance is not defined"))
		return
	}

	acc, err := a.store.UpdateAccount(context.Background(), repository.UpdateAccountParams{
		ID:      req.ID,
		Balance: *req.Balance,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusOK, genOkBody(acc))
}

type deleteAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func deleteAccount(ctx *gin.Context) {
	a := GetApi()
	req := &deleteAccountRequest{}

	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	if err := a.store.DeleteAccount(context.Background(), req.Id); err != nil {
		ctx.JSON(http.StatusBadRequest, genFailBody(err))
		return
	}

	ctx.JSON(http.StatusOK, genOkBody(struct{ Id int64 }{Id: req.Id}))
}
