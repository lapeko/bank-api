package db

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	want := CreateUserParams{
		FullName:       "Full Name",
		Email:          "random@mail.com",
		HashedPassword: strconv.Itoa(int(rand.Int63())),
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
