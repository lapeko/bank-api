package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
)

type authClientErrorMessage string

const (
	emailDuplicate      authClientErrorMessage = "email duplicated"
	emailNotFound       authClientErrorMessage = "user with given email does not exist"
	wrongPassword       authClientErrorMessage = "wrong password"
	invalidRefreshToken authClientErrorMessage = "provided refresh token is invalid"
)

type authClientError struct {
	message authClientErrorMessage
}

func (e *authClientError) Error() string {
	return string(e.message)
}

func genAccessToken(userId int64) (string, error) {
	now := time.Now()
	claims := utils.JWTUserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(utils.GetJwtKey())
}

func genRefreshToken(userId int64) (string, error) {
	now := time.Now()
	claims := utils.JWTUserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(14 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(utils.GetJwtKey())
}

func genTokens(userId int64) (tokens, error) {
	var tkns tokens
	accessToken, err := genAccessToken(userId)
	if err != nil {
		return tkns, err
	}
	refreshToken, err := genRefreshToken(userId)
	if err != nil {
		return tkns, err
	}
	tkns.AccessToken = accessToken
	tkns.RefreshToken = refreshToken
	return tkns, nil
}
