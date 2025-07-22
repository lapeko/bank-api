package auth

import (
	"context"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type authService struct {
	store db.Store
}

func (s *authService) createUser(ctx context.Context, args db.CreateUserParams) (user utils.UserResponse, err error) {
	dbUser, err := s.store.CreateUser(ctx, args)
	user = utils.DbUserToUserResponse(dbUser)
	return
}
