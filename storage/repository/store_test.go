package repository

import (
	"context"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	truncateTables()
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
		require.NotNil(t, transferRes.Transfer)
		require.NotNil(t, transferRes.EntryTo)
		require.NotNil(t, transferRes.EntryFrom)

		require.NotEmpty(t, transferRes.Transfer.ID)
		require.Equal(t, transferAmount, transferRes.Transfer.Amount)
		require.Equal(t, accFrom.ID, transferRes.Transfer.AccountFrom)
		require.Equal(t, accTo.ID, transferRes.Transfer.AccountTo)
		require.NotEmpty(t, transferRes.Transfer.CreatedAt)

		require.NotEmpty(t, transferRes.EntryFrom.ID)
		require.Equal(t, transferAmount, -transferRes.EntryFrom.Amount)
		require.Equal(t, accFrom.ID, transferRes.EntryFrom.AccountID)
		require.NotNil(t, transferRes.EntryFrom.CreatedAt)

		require.NotEmpty(t, transferRes.EntryTo.ID)
		require.Equal(t, transferAmount, transferRes.EntryTo.Amount)
		require.Equal(t, accTo.ID, transferRes.EntryTo.AccountID)
		require.NotNil(t, transferRes.EntryTo.CreatedAt)

		delta := transferAmount * (int64(i) + 1)
		require.Equal(t, accFrom.Balance-delta, transferRes.AccountFrom.Balance)
		require.Equal(t, accTo.Balance+delta, transferRes.AccountTo.Balance)
	}
}

func TestTransferTxDeadLock(t *testing.T) {
	truncateTables()
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
