package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/config"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	_ "github.com/lib/pq"
	"log"
)

type Api struct {
	router *gin.RouterGroup
	store  *repository.Store
	config *config.ApiConfig
}

func New(config *config.ApiConfig) *Api {
	a := &Api{}

	a.config = config

	db, err := sql.Open(a.config.DbDrives, a.config.DbAddress)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("DB successfully connected")

	a.store = repository.NewStore(db)

	log.Println("Store created")

	return a
}

func (a *Api) Start() {
	r := gin.Default()

	r.POST("/accounts", a.createAccount)
	r.GET("/accounts", a.getAccounts)
	r.GET("/accounts/:id", a.getAccount)
	r.PUT("/accounts", a.updateAccount)
	r.DELETE("/accounts/:id", a.deleteAccount)

	log.Fatalln(r.Run(a.config.ApiAddress))
}
