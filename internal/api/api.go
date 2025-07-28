package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/health"
	v1 "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/domains/v1"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type Api interface {
	Start(int) error
}

type api struct {
	router *gin.Engine
}

func New(conn *pgxpool.Pool) Api {
	router := gin.Default()
	utils.RegisterValidators()
	store := db.NewStore(conn)

	v1.Register("/v1", router, store)
	health.Register("/health", router)

	return &api{router: router}
}

func (a *api) Start(port int) error {
	return a.router.Run(fmt.Sprintf(":%d", port))
}
