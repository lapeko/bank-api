package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

type Api interface {
	Start(int) error
}

type api struct {
	router *gin.Engine
}

func New(conn *pgxpool.Pool) Api {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", utils.CurrencyValidator)
		v.RegisterValidation("fullname", utils.FullNameValidator)
		v.RegisterValidation("password", utils.PasswordValidator)
	}
	store := db.NewStore(conn)
	v1.Register("/v1", router, store)
	return &api{router: router}
}

func (a *api) Start(port int) error {
	return a.router.Run(fmt.Sprintf(":%d", port))
}
