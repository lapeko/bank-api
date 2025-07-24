package auth

import "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/v1/utils"

type tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required,fullname"`
	signinRequest
}

type createUserResponse struct {
	User utils.UserWithoutPassword `json:"user"`
	tokens
}

type signinRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
