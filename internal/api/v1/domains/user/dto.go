package user

import "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"

type listUsersRequest struct {
	Page int32 `form:"page" binding:"required,gte=1"`
	Size int32 `form:"size" binding:"required,gte=5,lte=20"`
}

type listUsersResponse struct {
	Users      []utils.UserWithoutPassword `json:"users"`
	TotalCount int64                       `json:"total_count"`
}

type updateUserEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type updateUserFullNameRequest struct {
	NewFullName string `json:"full_name" binding:"required,fullname"`
}

type updateUserPasswordRequest struct {
	NewPassword string `json:"password" binding:"required,password"`
}
