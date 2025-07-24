package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/domains/account"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/domains/auth"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/domains/entry"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/domains/transfer"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/domains/user"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

func Register(path string, router *gin.Engine, store db.Store) {
	g := router.Group(path)

	account.Register("/accounts", g, store)
	auth.Register("/auth", g, store)
	entry.Register("/entries", g, store)
	transfer.Register("/transfers", g, store)
	user.Register("/users", g, store)
}
