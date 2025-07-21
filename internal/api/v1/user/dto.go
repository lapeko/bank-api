package user

import "github.com/jackc/pgx/v5/pgtype"

type userResponse struct {
	ID        int64              `json:"id"`
	FullName  string             `json:"fullName"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type createUserRequest struct {
	FullName string `json:"fullName" binding:"required,fullname"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type listUsersRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listUsersResponse struct {
	Users      []userResponse `json:"users"`
	TotalCount int64          `json:"totalCount"`
}

type userUriIdRequest struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}

type updateUserEmailRequest struct {
	NewEmail string `json:"newEmail" binding:"required,email"`
}

type updateUserFullNameRequest struct {
	NewFullName string `json:"newFullName" binding:"required,fullname"`
}

type updateUserPasswordRequest struct {
	NewPassword string `json:"password" binding:"required,password"`
}
