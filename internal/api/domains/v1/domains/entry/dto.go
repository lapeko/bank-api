package entry

import (
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type listEntriesRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listEntriesResponse struct {
	Entries    []db.Entry `json:"entries"`
	TotalCount int64      `json:"total_count"`
}
