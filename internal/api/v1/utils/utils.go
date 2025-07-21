package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendError(ctx *gin.Context, err error) {
	SendErrorWithStatusCode(ctx, err, http.StatusBadRequest)
}

func SendErrorWithStatusCode(ctx *gin.Context, err error, statusCode int) {
	ctx.JSON(statusCode, genResponsePayload(nil, err.Error()))
}

func SendSuccess(ctx *gin.Context, payload interface{}) {
	SendSuccessWithStatusCode(ctx, payload, http.StatusOK)
}

func SendSuccessWithStatusCode(ctx *gin.Context, payload interface{}, statusCode int) {
	ctx.JSON(statusCode, genResponsePayload(payload, ""))
}

func genResponsePayload(payload interface{}, error string) gin.H {
	return gin.H{"error": error, "data": payload, "ok": error == ""}
}
