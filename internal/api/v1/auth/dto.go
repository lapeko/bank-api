package auth

type createUserRequest struct {
	FullName string `json:"fullName" binding:"required,fullname"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}
