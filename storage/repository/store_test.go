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
	require.Equal(t, err, errors.New("insufficient funds"))
}

func TestTransferTx(t *testing.T) {
	deleteAccounts(t)
	accFrom := createTestAccountWithBalance(t, random.Int64(100, 1000))
	accTo := createTestAccount(t)

	transferAmount := int64(10)

	const N = 5

	transferChan := make(chan *TransferTxResult)
	errorChan := make(chan error)

	for i := 0; i < N; i++ {
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

	for i := 0; i < N; i++ {
		err := <-errorChan
		require.Nil(t, err)

		transferRes := <-transferChan

		require.NotNil(t, transferRes)
		require.NotNil(t, transferRes.transfer)
		require.NotNil(t, transferRes.entryTo)
		require.NotNil(t, transferRes.entryFrom)

		require.NotEmpty(t, transferRes.transfer.ID)
		require.Equal(t, transferRes.transfer.Amount, transferAmount)
		require.Equal(t, transferRes.transfer.AccountFrom, accFrom.ID)
		require.Equal(t, transferRes.transfer.AccountTo, accTo.ID)
		require.NotEmpty(t, transferRes.transfer.CreatedAt)

		require.NotEmpty(t, transferRes.entryFrom.ID)
		require.Equal(t, -transferRes.entryFrom.Amount, transferAmount)
		require.Equal(t, transferRes.entryFrom.AccountID, accFrom.ID)
		require.NotNil(t, transferRes.entryFrom.CreatedAt)

		require.NotEmpty(t, transferRes.entryTo.ID)
		require.Equal(t, transferRes.entryTo.Amount, transferAmount)
		require.Equal(t, transferRes.entryTo.AccountID, accTo.ID)
		require.NotNil(t, transferRes.entryTo.CreatedAt)
	}
}
