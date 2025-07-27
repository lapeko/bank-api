package db

import (
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1, acc2 Account) Transfer {
	randAmount := int64(utils.GenRandIntInRange(-1e6, 1e6))

	got, err := testStore.CreateTransfer(ctx, CreateTransferParams{
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

func TestGetTotalTransfersCount(t *testing.T) {
	defer cleanTestStore(t)

	acc1 := createRandomAccount(t, createRandomUser(t))
	acc2 := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 3; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	count, err := testStore.GetTotalTransfersCount(ctx)
	require.NoError(t, err)
	require.True(t, count >= 3)
}

func TestGetTotalTransfersCountByAccount(t *testing.T) {
	defer cleanTestStore(t)

	acc1 := createRandomAccount(t, createRandomUser(t))
	acc2 := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 4; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	count, err := testStore.GetTotalTransfersCountByAccount(ctx, acc1.ID)
	require.NoError(t, err)
	require.Equal(t, int64(4), count)
}

func TestListTransfers(t *testing.T) {
	defer cleanTestStore(t)

	acc1 := createRandomAccount(t, createRandomUser(t))
	acc2 := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	transfers, err := testStore.ListTransfers(ctx, ListTransfersParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListTransfersByAccount(t *testing.T) {
	defer cleanTestStore(t)

	acc1 := createRandomAccount(t, createRandomUser(t))
	acc2 := createRandomAccount(t, createRandomUser(t))

	for i := 0; i < 7; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	transfers, err := testStore.ListTransfersByAccount(ctx, ListTransfersByAccountParams{
		AccountID: acc1.ID,
		Limit:     5,
		Offset:    0,
	})
	require.NoError(t, err)
	require.True(t, len(transfers) <= 5)
	for _, transfer := range transfers {
		require.True(t, transfer.AccountFrom == acc1.ID || transfer.AccountTo == acc1.ID)
	}
}
