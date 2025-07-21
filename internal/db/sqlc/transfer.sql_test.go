package db

import (
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1, acc2 Account) Transfer {
	randAmount := int64(utils.RandIntInRange(-1e6, 1e6))

	got, err := testStore.GetQueries().CreateTransfer(ctx, CreateTransferParams{
		AccountFrom: acc1.ID,
		AccountTo:   acc2.ID,
		Amount:      randAmount,
	})

	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.NotEmpty(t, got.ID)
	require.Equal(t, acc1.ID, got.AccountFrom)
	require.Equal(t, acc2.ID, got.AccountTo)
	require.Equal(t, got.Amount, randAmount)
	require.NotEmpty(t, got.CreatedAt)

	return got
}

func TestCreateTransfer(t *testing.T) {
	defer cleanTestStore(t)

	acc1 := createRandomAccount(t, createRandomUser(t))
	acc2 := createRandomAccount(t, createRandomUser(t))
	createRandomTransfer(t, acc1, acc2)
}
