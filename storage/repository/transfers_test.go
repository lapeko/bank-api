package repository

import (
	"context"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T, accountFromID int64, accountToID int64) Transfer {
	params := CreateTransferParams{
		AccountFrom: accountFromID,
		AccountTo:   accountToID,
		Amount:      random.Int64(1, 1000),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.AccountFrom, transfer.AccountFrom)
	require.Equal(t, params.AccountTo, transfer.AccountTo)
	require.Equal(t, params.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	params := CreateTransferParams{
		AccountFrom: account1.ID,
		AccountTo:   account2.ID,
		Amount:      random.Int64(1, 1000),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.AccountFrom, transfer.AccountFrom)
	require.Equal(t, params.AccountTo, transfer.AccountTo)
	require.Equal(t, params.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	transfer := createTestTransfer(t, account1.ID, account2.ID)

	storedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, storedTransfer)

	require.Equal(t, transfer.ID, storedTransfer.ID)
	require.Equal(t, transfer.AccountFrom, storedTransfer.AccountFrom)
	require.Equal(t, transfer.AccountTo, storedTransfer.AccountTo)
	require.Equal(t, transfer.Amount, storedTransfer.Amount)
}

func TestListTransfers(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	for i := 0; i < 5; i++ {
		createTestTransfer(t, account1.ID, account2.ID)
	}

	params := ListTransfersParams{
		Limit:  5,
		Offset: 0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
}

func TestListTransfersByReceiver(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	for i := 0; i < 5; i++ {
		createTestTransfer(t, account1.ID, account2.ID)
	}

	params := ListTransfersByReceiverParams{
		AccountTo: account2.ID,
		Limit:     5,
		Offset:    0,
	}

	transfers, err := testQueries.ListTransfersByReceiver(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.Equal(t, account2.ID, transfer.AccountTo)
	}
}

func TestListTransfersBySender(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	for i := 0; i < 5; i++ {
		createTestTransfer(t, account1.ID, account2.ID)
	}

	params := ListTransfersBySenderParams{
		AccountFrom: account1.ID,
		Limit:       5,
		Offset:      0,
	}

	transfers, err := testQueries.ListTransfersBySender(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.Equal(t, account1.ID, transfer.AccountFrom)
	}
}

func TestListTransfersBySenderAndReceiver(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	for i := 0; i < 5; i++ {
		createTestTransfer(t, account1.ID, account2.ID)
	}

	params := ListTransfersBySenderAndReceiverParams{
		AccountFrom: account1.ID,
		AccountTo:   account2.ID,
		Limit:       5,
		Offset:      0,
	}

	transfers, err := testQueries.ListTransfersBySenderAndReceiver(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.Equal(t, account1.ID, transfer.AccountFrom)
		require.Equal(t, account2.ID, transfer.AccountTo)
	}
}
