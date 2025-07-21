package account

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

type accountWithUserInfo struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"userId"`
	FullName  string             `json:"fullName"`
	Email     string             `json:"email"`
	Currency  utils.Currency     `json:"currency"`
	Balance   int64              `json:"balance"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type createAccountRequest struct {
	UserID   int64          `json:"userId" binding:"required,gte=1"`
	Currency utils.Currency `json:"currency" binding:"required,currency"`
}

type listAccountsRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listAccountsResponse struct {
	Accounts   []accountWithUserInfo `json:"accounts"`
	TotalCount int64                 `json:"totalCount"`
}

type getAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}

type getAccountByIdResponse struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"userId"`
	FullName  string             `json:"fullName"`
	Email     string             `json:"email"`
	Currency  utils.Currency     `json:"currency"`
	Balance   int64              `json:"balance"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type deleteAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}
