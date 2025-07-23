package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type UserWithoutPassword struct {
	ID        int64              `json:"id"`
	FullName  string             `json:"full_name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func CutUserPassword(u db.User) UserWithoutPassword {
	return UserWithoutPassword{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
