package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	_ "github.com/lib/pq"
	"log"
)

type Api struct {
	router *gin.RouterGroup
	store  *repository.Store
}

func New(postgresAddr string) *Api {
	api := &Api{}

	db, err := sql.Open("postgres", postgresAddr)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("DB successfully connected")

	api.store = repository.NewStore(db)

	log.Println("Store created")

	return api
}

func (a *Api) Start(apiAddr string) {
	r := gin.Default()

	r.POST("/accounts", a.createAccount)
	r.GET("/accounts", a.getAccounts)

	log.Fatalln(r.Run(apiAddr))
}
