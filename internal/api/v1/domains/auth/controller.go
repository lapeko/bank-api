package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	rootUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

var service authService

func Register(path string, router *gin.RouterGroup, store db.Store) {
	service = authService{store: store}
	g := router.Group(path)

	g.POST("/signup", signupHandler)
	g.POST("/signin", signinHandler)
	g.POST("/refresh", refreshHandler)
}

func signupHandler(ctx *gin.Context) {
	var usr createUserRequest
	if err := ctx.ShouldBind(&usr); err != nil {
		utils.SendError(ctx, err)
		return
	}
	hash, err := rootUtils.HashPassword(usr.Password)
	if err != nil {
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	params := db.CreateUserParams{FullName: usr.FullName, Email: usr.Email, HashedPassword: hash}
	res, err := service.createUser(ctx, params)
	if err != nil {
		// TODO handle email taken
		utils.SendErrorWithStatusCode(ctx, err, http.StatusInternalServerError)
		return
	}
	utils.SendSuccessWithStatusCode(ctx, res, http.StatusCreated)
}

func signinHandler(ctx *gin.Context) {

}

func refreshHandler(ctx *gin.Context) {

}
