package transfer

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

type uriIdRequest struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}

type transferResponse struct {
	ID          int64              `json:"id"`
	AccountFrom int64              `json:"accountFrom"`
	AccountTo   int64              `json:"accounTo"`
	Amount      int64              `json:"amount"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
}

type transferRequest struct {
	AccountFrom int64 `json:"accountFrom" binding:"required,gte=1"`
	AccountTo   int64 `json:"accountTo" binding:"required,gte=1"`
	Amount      int64 `json:"amount" binding:"required,gte=1"`
}

type externalTransferRequest struct {
	Account int64 `json:"accountFrom" binding:"required,gte=1"`
	Amount  int64 `json:"amount" binding:"required"`
}

type accountResponse struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"userId"`
	Currency  utils.Currency     `json:"currency"`
	Balance   int64              `json:"balance"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type listTransfersRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listTransfersResponse struct {
	Transfers  []transferResponse `json:"transfers"`
	TotalCount int64              `json:"totalCount"`
}
