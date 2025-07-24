package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UriId struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}

type UserWithoutPassword struct {
	ID        int64              `json:"id"`
	FullName  string             `json:"full_name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
