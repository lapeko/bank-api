package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/account"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

func Register(path string, router *gin.Engine, store db.Store) {
	g := router.Group(path)
	account.Register("/accounts", g, store)
}
