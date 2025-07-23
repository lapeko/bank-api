package user

import (
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/v1/utils"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

func dbUsersToUserResponses(src []db.User) []utils.UserWithoutPassword {
	users := make([]utils.UserWithoutPassword, len(src))
	for idx, dbUser := range src {
		users[idx] = utils.CutUserPassword(dbUser)
	}
	return users
}
