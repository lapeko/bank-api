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

		delta := transferAmount * (int64(i) + 1)
		require.Equal(t, accFrom.Balance-delta, transferRes.accountFrom.Balance)
		require.Equal(t, accTo.Balance+delta, transferRes.accountTo.Balance)
	}
}

func TestTransferTxDeadLock(t *testing.T) {
	deleteAccounts(t)
	acc1 := createTestAccountWithBalance(t, random.Int64(100, 1000))
	acc2 := createTestAccountWithBalance(t, random.Int64(100, 1000))

	transferAmount := int64(10)

	const n = 10

	errorChan := make(chan error)

	for i := 0; i < n; i++ {
		AccountFrom := acc1.ID
		AccountTo := acc2.ID

		if i%2 == 0 {
			AccountFrom = acc2.ID
			AccountTo = acc1.ID
		}

		go func() {
			_, err := NewStore(testDb).TransferTX(context.Background(), CreateTransferParams{
				AccountFrom: AccountFrom,
				AccountTo:   AccountTo,
				Amount:      transferAmount,
			})
			errorChan <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errorChan
		require.Nil(t, err)
	}

	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Nil(t, err)
	require.NotNil(t, updatedAcc1)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.Nil(t, err)
	require.NotNil(t, updatedAcc2)

	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
