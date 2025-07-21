package user

import (
	"context"

	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"golang.org/x/sync/errgroup"
)

type userService struct {
	store db.Store
}

func (s *userService) createUser(ctx context.Context, args db.CreateUserParams) (user userResponse, err error) {
	dbUser, err := s.store.CreateUser(ctx, args)
	user = dbUserToUserResponse(dbUser)
	return
}

func (s *userService) listUsers(ctx context.Context, args listUsersRequest) (res listUsersResponse, err error) {
	params := db.ListUsersParams{
		Limit:  args.Size,
		Offset: (args.Page - 1) * args.Size,
	}

	var users []db.User
	var count int64

	g := errgroup.Group{}
	g.Go(func() (e error) {
		users, e = s.store.ListUsers(ctx, params)
		return
	})
	g.Go(func() (e error) {
		count, e = s.store.GetTotalUsersCount(ctx)
		return
	})
	if err = g.Wait(); err != nil {
		return res, err
	}

	return listUsersResponse{
		Users:      dbUsersToUserResponses(users),
		TotalCount: count,
	}, nil
}

func (s *userService) getUserById(ctx context.Context, id int64) (user userResponse, err error) {
	dbUser, err := s.store.GetUserById(ctx, id)
	user = dbUserToUserResponse(dbUser)
	return
}

func (s *userService) updateUserEmail(ctx context.Context, id int64, newEmail string) (userResponse, error) {
	dbUser, err := s.store.UpdateUserEmail(ctx, db.UpdateUserEmailParams{ID: id, Email: newEmail})
	return dbUserToUserResponse(dbUser), err
}

func (s *userService) updateUserFullName(ctx context.Context, id int64, newFullName string) (userResponse, error) {
	dbUser, err := s.store.UpdateUserFullName(ctx, db.UpdateUserFullNameParams{ID: id, FullName: newFullName})
	return dbUserToUserResponse(dbUser), err
}

func (s *userService) updateUserPassword(ctx context.Context, id int64, hashedPassword string) (userResponse, error) {
	dbUser, err := s.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{ID: id, HashedPassword: hashedPassword})
	return dbUserToUserResponse(dbUser), err
}

func (s *userService) deleteUser(ctx context.Context, id int64) (err error) {
	_, err = s.store.DeleteUser(ctx, id)
	return
}
