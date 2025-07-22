package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type UserResponse struct {
	ID        int64              `json:"id"`
	FullName  string             `json:"fullName"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

func DbUserToUserResponse(src db.User) UserResponse {
	return UserResponse{
		ID:        src.ID,
		FullName:  src.FullName,
		Email:     src.Email,
		CreatedAt: src.CreatedAt,
	}
}
