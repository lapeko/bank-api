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
	router *gin.Engine
	store  repository.Store
	config *config.ApiConfig
}

var api *Api

func GetApi() *Api {
	if api == nil {
		api = &Api{}
	}
	return api
}

func (a *Api) ConnectStore(config *config.ApiConfig) {
	a.config = config

	db, err := sql.Open(a.config.DbDrives, a.config.DbAddress)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("DB successfully connected")

	a.store = repository.NewStore(db)

	log.Println("Store created")
}

func (a *Api) SetUpRoutes() {
	a.router = gin.Default()

	setUpAccounts(a.router)
}

func (a *Api) Start() {
	log.Fatalln(a.router.Run(a.config.ApiAddress))
}
