package transfer

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

var service transferService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = transferService{store: store}
	g := router.Group(path)
	g.POST("/", transferHandler)
	g.POST("/external", externalTransferHandler)
	g.GET("/", listTransfersHandler)
	g.GET("/account/:id", listTransfersByAccountHandler)
}

func transferHandler(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := service.transferMoney(ctx, req); err != nil {
		var transferClientError *db.TransferClientError
		if errors.As(err, &transferClientError) {
			utils.SendError(ctx, err)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, struct{}{})
}

func externalTransferHandler(ctx *gin.Context) {

}

func listTransfersHandler(ctx *gin.Context) {

}

func listTransfersByAccountHandler(ctx *gin.Context) {

}
