package auth

import (
	"context"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	internalUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

type authService struct {
	store db.Store
}

type authClientErrorMessage string

const (
	emailNotFound       authClientErrorMessage = "user with given email does not exist"
	wrongPassword       authClientErrorMessage = "wrong password"
	refreshTokenExpired authClientErrorMessage = "refresh token expired"
)

type authClientError struct {
	message authClientErrorMessage
}

func (e *authClientError) Error() string {
	return string(e.message)
}

func (s *authService) createUser(ctx context.Context, args db.CreateUserParams) (res createUserResponse, err error) {
	user, err := s.store.CreateUser(ctx, args)
	if err != nil {
		return
	}
	res.User = utils.CutUserPassword(user)

	tkns, err := genTokens(user.ID)
	if err != nil {
		return
	}

	res.tokens = tkns
	return
}

func (s *authService) signIn(ctx context.Context, args signInRequest) (tkns tokens, err error) {
	user, err := s.store.GetUserByEmail(ctx, args.Email)
	if err != nil {
		return tkns, &authClientError{emailNotFound}
	}
	if ok := internalUtils.CompareHashAndPassword(user.HashedPassword, args.Password); !ok {
		return tkns, &authClientError{wrongPassword}
	}
	tkns, err = genTokens(user.ID)
	return
}
