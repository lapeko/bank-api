package user

import "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"

type listUsersRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listUsersResponse struct {
	Users      []utils.UserResponse `json:"users"`
	TotalCount int64                `json:"totalCount"`
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
