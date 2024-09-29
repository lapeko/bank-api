package repository

import (
	"context"
	"errors"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferInsufficientFunds(t *testing.T) {
	deleteAccounts(t)
	accTo := createTestAccount(t)
	accFrom := createTestAccountWithBalance(t, random.Int64(-1000, -100))

	store := NewStore(testDb)

	transfer, err := store.TransferTX(context.Background(), CreateTransferParams{
		AccountFrom: accFrom.ID,
		AccountTo:   accTo.ID,
		Amount:      random.Int64(1, 1000),
	})

	require.Nil(t, transfer)
	require.Equal(t, errors.New("insufficient funds"), err)
}

func TestTransferTx(t *testing.T) {
	deleteAccounts(t)
	accFrom := createTestAccountWithBalance(t, random.Int64(100, 1000))
	accTo := createTestAccount(t)

	transferAmount := int64(10)

	const n = 5

	transferChan := make(chan *TransferTxResult)
	errorChan := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			transferRes, err := NewStore(testDb).TransferTX(context.Background(), CreateTransferParams{
				AccountFrom: accFrom.ID,
				AccountTo:   accTo.ID,
				Amount:      transferAmount,
			})
			errorChan <- err
			transferChan <- transferRes
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errorChan
		require.Nil(t, err)

		transferRes := <-transferChan

		require.NotNil(t, transferRes)
		require.NotNil(t, transferRes.transfer)
		require.NotNil(t, transferRes.entryTo)
		require.NotNil(t, transferRes.entryFrom)

		require.NotEmpty(t, transferRes.transfer.ID)
		require.Equal(t, transferAmount, transferRes.transfer.Amount)
		require.Equal(t, accFrom.ID, transferRes.transfer.AccountFrom)
		require.Equal(t, accTo.ID, transferRes.transfer.AccountTo)
		require.NotEmpty(t, transferRes.transfer.CreatedAt)

		require.NotEmpty(t, transferRes.entryFrom.ID)
		require.Equal(t, transferAmount, -transferRes.entryFrom.Amount)
		require.Equal(t, accFrom.ID, transferRes.entryFrom.AccountID)
		require.NotNil(t, transferRes.entryFrom.CreatedAt)

		require.NotEmpty(t, transferRes.entryTo.ID)
		require.Equal(t, transferAmount, transferRes.entryTo.Amount)
		require.Equal(t, accTo.ID, transferRes.entryTo.AccountID)
		require.NotNil(t, transferRes.entryTo.CreatedAt)

		updatedAccFrom, err := testQueries.GetAccount(context.Background(), accFrom.ID)
		require.Nil(t, err)
		require.NotNil(t, updatedAccFrom)

		updatedAccTo, err := testQueries.GetAccount(context.Background(), accTo.ID)
		require.Nil(t, err)
		require.NotNil(t, updatedAccTo)

		delta := transferAmount * (int64(i) + 1)
		require.Equal(t, accFrom.Balance-delta, updatedAccFrom.Balance)
		require.Equal(t, accTo.Balance+delta, updatedAccTo.Balance)
	}
}
