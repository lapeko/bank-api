package db

import (
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	randAmount := int64(utils.RandIntInRange(-1e6, 1e6))

	got, err := testStore.CreateEntry(ctx, CreateEntryParams{
		AccountID: acc.ID,
		Amount:    randAmount,
	})

	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.NotEmpty(t, got.ID)
	require.Equal(t, acc.ID, got.AccountID)
	require.Equal(t, got.Amount, randAmount)
	require.NotEmpty(t, got.CreatedAt)

	return got
}

func TestCreateEntry(t *testing.T) {
	defer cleanTestStore(t)

	createRandomEntry(t, createRandomAccount(t, createRandomUser(t)))
}
