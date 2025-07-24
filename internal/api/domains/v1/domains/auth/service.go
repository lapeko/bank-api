package auth

import (
	"context"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/v1/utils"
	apiUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type authService struct {
	store db.Store
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

func (s *authService) signIn(ctx context.Context, args signinRequest) (tkns tokens, err error) {
	user, err := s.store.GetUserByEmail(ctx, args.Email)
	if err != nil {
		return tkns, &authClientError{emailNotFound}
	}
	if ok := apiUtils.CompareHashAndPassword(user.HashedPassword, args.Password); !ok {
		return tkns, &authClientError{wrongPassword}
	}
	tkns, err = genTokens(user.ID)
	return
}

func (s *authService) refreshToken(ctx context.Context, args refreshTokenRequest) (rtknRes refreshTokenResponse, err error) {
	claims, ok := apiUtils.ParseJwtToken(args.RefreshToken)
	if !ok {
		return rtknRes, &authClientError{invalidRefreshToken}
	}
	rtknRes.AccessToken, err = genAccessToken(claims.UserId)
	return
}
