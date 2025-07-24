package middleware

import "github.com/gin-gonic/gin"

const UserInfoKey = "user_info"

type UserInfo struct {
	Id int64
}

func GetUserInfo(ctx *gin.Context) (userInfo UserInfo, ok bool) {
	info, ok := ctx.Get(UserInfoKey)
	if !ok {
		return
	}
	userInfo, ok = info.(UserInfo)
	return
}

func MustGetUserInfo(ctx *gin.Context) UserInfo {
	info, ok := GetUserInfo(ctx)
	if !ok {
		panic("unable to get user info from context")
	}
	return info
}
