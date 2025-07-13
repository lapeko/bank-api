package db

import (
	"context"
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/db/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	want := CreateUserParams{
		FullName:       utils.GenRandFullName(),
		Email:          utils.GetRandEmail(),
		HashedPassword: utils.GenRandHashedPassword(),
	}
	got, err := sqlcQueries.CreateUser(context.TODO(), want)

	require.NoError(t, err)
	require.NotZero(t, got.ID)
	require.Equal(t, got.FullName, want.FullName)
	require.Equal(t, got.Email, want.Email)
	require.Equal(t, got.HashedPassword, want.HashedPassword)
	require.NotZero(t, got.CreatedAt)

	return &got
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
