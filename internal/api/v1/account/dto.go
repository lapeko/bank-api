package account

import (
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

type createAccountRequest struct {
	UserID   int64          `json:"user_id" binding:"required,gte=1"`
	Currency utils.Currency `json:"currency" binding:"required,currency"`
}

type listAccountsRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listAccountsResponse struct {
	Accounts   []db.ListAccountsRow `json:"accounts"`
	TotalCount int64                `json:"total_count"`
}
