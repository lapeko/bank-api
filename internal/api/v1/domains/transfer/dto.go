package transfer

import (
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type transferRequest struct {
	AccountFrom int64 `json:"account_from" binding:"required,gte=1"`
	AccountTo   int64 `json:"account_to" binding:"required,gte=1"`
	Amount      int64 `json:"amount" binding:"required,gte=1"`
}

type externalTransferRequest struct {
	Account int64 `json:"account_from" binding:"required,gte=1"`
	Amount  int64 `json:"amount" binding:"required"`
}

type listTransfersRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listTransfersResponse struct {
	Transfers  []db.Transfer `json:"transfers"`
	TotalCount int64         `json:"total_count"`
}
