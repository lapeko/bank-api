package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUser(t *testing.T) {
	truncateTables()
	user := createTestUser(t)

	gotUser, err := testQueries.GetUser(context.Background(), user.ID)

	require.Nil(t, err)
	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.CreatedAt, gotUser.CreatedAt)
	require.Equal(t, user.HashedPassword, gotUser.HashedPassword)
	require.Equal(t, user.PasswordChangesAt, gotUser.PasswordChangesAt)
}

func TestGetNotExistUser(t *testing.T) {
	truncateTables()

	user, err := testQueries.GetUser(context.Background(), 1)

	require.Zero(t, user)
	require.Equal(t, err, sql.ErrNoRows)
}

func createTestUser(t *testing.T) *User {
	params := CreateUserParams{
		FullName:       random.String(10),
		Email:          fmt.Sprintf("%s@email.com", random.String(6)),
		HashedPassword: "secret",
	}

	user, err := testQueries.CreateUser(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.FullName, params.FullName)
	require.Equal(t, user.Email, user.Email)
	require.Equal(t, user.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return &user
}
