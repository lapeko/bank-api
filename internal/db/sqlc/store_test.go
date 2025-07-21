package db

import (
	"errors"
	"sync"
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const zero = int64(0)

func createTwoAccountsWithBalances(t *testing.T, balance1, balance2 int64, currency utils.Currency) (Account, Account) {
	acc1 := createAccountWithParams(t, CreateAccountParams{UserID: createRandomUser(t).ID, Currency: currency})
	acc2 := createAccountWithParams(t, CreateAccountParams{UserID: createRandomUser(t).ID, Currency: currency})

	acc1, err := testStore.OffsetBalance(ctx, OffsetBalanceParams{ID: acc1.ID, Delta: balance1})
	require.NoError(t, err)
	acc2, err = testStore.OffsetBalance(ctx, OffsetBalanceParams{ID: acc2.ID, Delta: balance2})
	require.NoError(t, err)

	require.Equal(t, balance1, acc1.Balance)
	require.Equal(t, balance2, acc2.Balance)

	return acc1, acc2
}

func queryTwoAccountsById(t *testing.T, acc1Id, acc2Id int64) (GetAccountByIdRow, GetAccountByIdRow) {
	acc1, err := testStore.GetAccountById(ctx, acc1Id)
	require.NoError(t, err)
	acc2, err := testStore.GetAccountById(ctx, acc2Id)
	require.NoError(t, err)
	return acc1, acc2
}

func TestTransferMoney(t *testing.T) {
	defer cleanTestStore(t)

	transferAmmount := utils.RandInt64InRange(1, 1e4)
	acc1, acc2 := createTwoAccountsWithBalances(t, zero, transferAmmount, utils.GenRandCurrency())

	err := testStore.TransferMoney(ctx, acc2.ID, acc1.ID, transferAmmount)
	require.NoError(t, err)

	extAcc1, extAcc2 := queryTwoAccountsById(t, acc1.ID, acc2.ID)
	require.Equal(t, extAcc1.Balance, transferAmmount)
	require.Equal(t, extAcc2.Balance, zero)
}

func TestTransferMoney_Concurrently(t *testing.T) {
	defer cleanTestStore(t)

	store := testStore

	accBalance := utils.RandInt64InRange(1e5, 1e6)
	acc2Balance := utils.RandInt64InRange(1e5, 1e6)
	transferLeft := utils.RandInt64InRange(1, 1e4)
	transferRight := utils.RandInt64InRange(1, 1e4)

	acc1, acc2 := createTwoAccountsWithBalances(t, accBalance, acc2Balance, utils.GenRandCurrency())

	iterations := 5
	size := iterations * 2
	errs := make(chan error, size)

	var wg sync.WaitGroup
	wg.Add(size)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			errs <- store.TransferMoney(ctx, acc2.ID, acc1.ID, transferLeft)
		}()
		go func() {
			defer wg.Done()
			errs <- store.TransferMoney(ctx, acc1.ID, acc2.ID, transferRight)
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	delta := (transferRight - transferLeft) * int64(iterations)
	extAcc1, extAcc2 := queryTwoAccountsById(t, acc1.ID, acc2.ID)
	require.Equal(t, extAcc1.Balance, accBalance-delta)
	require.Equal(t, extAcc2.Balance, acc2Balance+delta)
}

func TestTransferMoney_NotPositiveAmountError(t *testing.T) {
	defer cleanTestStore(t)

	acc1, acc2 := createTwoAccountsWithBalances(t, zero, zero, utils.GenRandCurrency())

	err := testStore.TransferMoney(ctx, acc1.ID, acc2.ID, zero)
	require.Error(t, err)
	require.Equal(t, err.Error(), "amount should be positive")

	extAcc1, extAcc2 := queryTwoAccountsById(t, acc1.ID, acc2.ID)
	require.Equal(t, extAcc1.Balance, zero)
	require.Equal(t, extAcc2.Balance, zero)
}

func TestTransferMoney_SameAccountError(t *testing.T) {
	defer cleanTestStore(t)

	acc := createAccountWithParams(t, CreateAccountParams{UserID: createRandomUser(t).ID, Currency: utils.GenRandCurrency()})

	err := testStore.TransferMoney(ctx, acc.ID, acc.ID, 1)
	require.Error(t, err)
	require.Equal(t, err.Error(), "transfer accounts should not be same")

	extAcc, err := testStore.GetAccountById(ctx, acc.ID)
	require.NoError(t, err)
	require.Equal(t, extAcc.Balance, zero)
}

func TestTransferMoney_DifferentCurrencyError(t *testing.T) {
	defer cleanTestStore(t)

	acc1 := createAccountWithParams(t, CreateAccountParams{UserID: createRandomUser(t).ID, Currency: utils.CurrencyEURO})
	acc2 := createAccountWithParams(t, CreateAccountParams{UserID: createRandomUser(t).ID, Currency: utils.CurrencyUSD})

	err := testStore.TransferMoney(ctx, acc1.ID, acc2.ID, 1)
	require.Error(t, err)
	require.Equal(t, err.Error(), "currency must be same for money transfer")

	extAcc1, extAcc2 := queryTwoAccountsById(t, acc1.ID, acc2.ID)
	require.Equal(t, extAcc1.Balance, zero)
	require.Equal(t, extAcc2.Balance, zero)
}

func TestTransferMoney_InsufficientFundsError(t *testing.T) {
	defer cleanTestStore(t)

	acc1, acc2 := createTwoAccountsWithBalances(t, zero, zero, utils.GenRandCurrency())

	err := testStore.TransferMoney(ctx, acc1.ID, acc2.ID, 1)
	require.Error(t, err)
	require.Equal(t, err.Error(), "insufficient funds")

	extAcc1, extAcc2 := queryTwoAccountsById(t, acc1.ID, acc2.ID)
	require.Equal(t, extAcc1.Balance, zero)
	require.Equal(t, extAcc2.Balance, zero)
}

func TestTransferMoney_GetAccountsByIdForUpdateError(t *testing.T) {
	defer cleanTestStore(t)

	queryErr := errors.New(utils.GenRandString(10))
	rbcErrMsg := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	txMock := new(dbConnMock)
	dbMock.On("Begin", ctx).Return(txMock, nil)
	txMock.On("Query", ctx, getAccountsByIdForUpdate, mock.Anything, mock.Anything).Return(nil, queryErr)
	txMock.On("Rollback", ctx).Return(rbcErrMsg)

	testStoreMock := NewStore(dbMock)

	err := testStoreMock.TransferMoney(ctx, 1, 2, int64(1))
	require.Error(t, err)
	require.ErrorIs(t, err, queryErr)
	require.ErrorIs(t, err, rbcErrMsg)
}

func TestTransferMoney_TXBeginError(t *testing.T) {
	defer cleanTestStore(t)

	txError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	dbMock.On("Begin", ctx).Return(new(dbConnMock), txError)

	testStoreMock := NewStore(dbMock)

	err := testStoreMock.TransferMoney(ctx, 1, 2, int64(1))
	require.Error(t, err)
	require.ErrorIs(t, err, txError)
}
