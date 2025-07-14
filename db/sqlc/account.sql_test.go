package db

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/db/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, user User) Account {
	want := CreateAccountParams{
		UserID:   user.ID,
		Currency: utils.GenRandCurrency(),
		Balance:  rand.Int64N(1e6),
	}

	fmt.Println()
	fmt.Println(want.Currency)

	got, err := testStore.GetQueries().CreateAccount(ctx, want)

	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.NotZero(t, got.ID)
	require.Equal(t, got.UserID, want.UserID)
	require.Equal(t, got.Currency, want.Currency)
	require.Equal(t, got.Balance, want.Balance)
	require.NotZero(t, got.CreatedAt)

	return got
}

func TestCreateAccount(t *testing.T) {
	defer cleanTestStore(t)

	createRandomAccount(t, createRandomUser(t))
}

func TestDeleteAccount(t *testing.T) {
	defer cleanTestStore(t)

	q := testStore.GetQueries()
	acc := createRandomAccount(t, createRandomUser(t))
	accById, err := q.GetAccountById(ctx, acc.ID)

	require.NoError(t, err)
	require.Equal(t, acc, accById)

	err = q.DeleteAccount(ctx, acc.ID)
	require.NoError(t, err)

	accById, err = q.GetAccountById(ctx, acc.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Empty(t, accById)
}

func TestListAccounts(t *testing.T) {
	defer cleanTestStore(t)
	user := createRandomUser(t)

	params := ListAccountsParams{Offset: 0, Limit: 2}
	q := testStore.GetQueries()
	got, err := q.ListAccounts(ctx, params)
	require.NoError(t, err)
	require.Empty(t, got)

	var want []Account
	for i := 0; i < 2; i++ {
		want = append(want, createRandomAccount(t, user))
	}
	got, err = q.ListAccounts(ctx, params)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestListAccounts_QueryError(t *testing.T) {
	params := ListAccountsParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(DBTXMock)
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, wantError)

	gotAccounts, gotError := New(dbMock).ListAccounts(ctx, params)
	require.Nil(t, gotAccounts)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListAccounts_ScanError(t *testing.T) {
	params := ListAccountsParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(DBTXMock)
	rowsMock := new(RowsMock)
	rowsMock.On("Next").Return(true)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotAccounts, gotError := New(dbMock).ListAccounts(ctx, params)
	require.Nil(t, gotAccounts)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Close")
}

func TestListAccounts_RowsError(t *testing.T) {
	params := ListAccountsParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(DBTXMock)
	rowsMock := new(RowsMock)
	rowsMock.On("Next").Once().Return(true)
	rowsMock.On("Next").Once().Return(false)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	rowsMock.On("Err").Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotUsers, gotError := New(dbMock).ListAccounts(ctx, params)
	require.Nil(t, gotUsers)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Err")
	rowsMock.AssertCalled(t, "Close")
}

func TestOffsetBalance(t *testing.T) {
	defer cleanTestStore(t)

	acc := createRandomAccount(t, createRandomUser(t))
	deltaBalance := int64(utils.RandIntInRange(-1e6, 1e6))
	wantBalance := acc.Balance + deltaBalance

	gotAcc, err := testStore.GetQueries().OffsetBalance(ctx, OffsetBalanceParams{
		ID:    acc.ID,
		Delta: int64(deltaBalance),
	})

	require.NoError(t, err)
	require.Equal(t, wantBalance, gotAcc.Balance)
}

func TestUpdateAccountBalance(t *testing.T) {
	defer cleanTestStore(t)

	acc := createRandomAccount(t, createRandomUser(t))
	wantBalance := int64(utils.RandIntInRange(-1e6, 1e6))

	gotAcc, err := testStore.GetQueries().UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:      acc.ID,
		Balance: int64(wantBalance),
	})

	require.NoError(t, err)
	require.Equal(t, wantBalance, gotAcc.Balance)
}
