package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/v1/utils"
)

var unauthorized = errors.New("not authorized")

func AuthGuard(ctx *gin.Context) {
	if info, ok := ctx.Get(UserInfoKey); ok {
		if _, ok = info.(UserInfo); ok {
			ctx.Next()
			return
		}
	}
	utils.SendErrorWithStatusCode(ctx, unauthorized, http.StatusUnauthorized)
}
