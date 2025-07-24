package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type userClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func genAccessToken(userId int64) (string, error) {
	now := time.Now()
	claims := userClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(jwtKey)
}

func genRefreshToken(userId int64) (string, error) {
	now := time.Now()
	claims := userClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(14 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(jwtKey)
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

func parseToken(tokenString string) (claims *userClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(t *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		return
	}
	if claims, ok := token.Claims.(*userClaims); ok {
		return claims, nil
	}
	return claims, errors.New("jwt claim parse failure")
}
