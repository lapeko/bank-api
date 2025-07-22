package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(path string, router *gin.RouterGroup) {
	g := router.Group(path)
	g.GET("/", healthHandler)
}

func healthHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}
