package entry

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

var service entryService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = entryService{store: store}
	g := router.Group(path)

	g.GET("/", listEntriesHandler)
	g.GET("/:id", getEntryByIdHandler)
	g.GET("/account/:id", listEntriesByAccountHandler)
}

func listEntriesHandler(ctx *gin.Context) {
	var req listEntriesRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendError(ctx, err)
		return
	}
	payload, err := service.listEntries(ctx, req)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, payload)
}

func getEntryByIdHandler(ctx *gin.Context) {
	var uriId utils.UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	entry, err := service.getEntryById(ctx, uriId.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendErrorWithStatusCode(ctx, err, http.StatusNotFound)
			return
		}
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, entry)
}

func listEntriesByAccountHandler(ctx *gin.Context) {
	var uriId utils.UriId
	var queryParams listEntriesRequest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		utils.SendError(ctx, err)
		return
	}
	if err := ctx.ShouldBind(&queryParams); err != nil {
		utils.SendError(ctx, err)
		return
	}
	list, err := service.listEntriesByAccount(ctx, queryParams, uriId.ID)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(ctx, list)
}
