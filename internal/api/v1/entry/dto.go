package entry

import "github.com/jackc/pgx/v5/pgtype"

type uriIdRequest struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}

type entryResponse struct {
	ID        int64              `json:"id"`
	AccountID int64              `json:"accountId"`
	Amount    int64              `json:"amount"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type listEntriesRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listEntriesResponse struct {
	Entries    []entryResponse `json:"entries"`
	TotalCount int64           `json:"totalCount"`
}
