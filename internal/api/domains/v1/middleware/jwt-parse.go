package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
)

func JwtParser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, ok := utils.ParseJwtToken(tokenString); ok {
			ctx.Set(UserInfoKey, UserInfo{Id: claims.UserId})
		}
	}

	ctx.Next()
}
