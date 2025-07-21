package user

import db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"

func dbUserToUserResponse(src db.User) userResponse {
	return userResponse{
		ID:        src.ID,
		FullName:  src.FullName,
		Email:     src.Email,
		CreatedAt: src.CreatedAt,
	}
}

func dbUsersToUserResponses(src []db.User) []userResponse {
	users := make([]userResponse, len(src))
	for idx, dbUser := range src {
		users[idx] = dbUserToUserResponse(dbUser)
	}
	return users
}
