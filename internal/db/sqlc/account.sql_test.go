package db

import (
	"errors"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, user User) Account {
	want := CreateAccountParams{
		UserID:   user.ID,
		Currency: utils.GenRandCurrency(),
	}

	return createAccountWithParams(t, want)
}

func createAccountWithParams(t *testing.T, want CreateAccountParams) Account {
	got, err := testStore.CreateAccount(ctx, want)

	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.NotZero(t, got.ID)
	require.Equal(t, got.UserID, want.UserID)
	require.Equal(t, got.Currency, want.Currency)
	require.Equal(t, got.Balance, zero)
	require.NotZero(t, got.CreatedAt)

	return got
}

func TestCreateAccount(t *testing.T) {
	defer cleanTestStore(t)

	createRandomAccount(t, createRandomUser(t))
}

func TestGetAccountByIdForUpdate(t *testing.T) {
	defer cleanTestStore(t)

	acc := createRandomAccount(t, createRandomUser(t))

	accById, err := testStore.GetAccountByIdForUpdate(ctx, acc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accById)
	require.Equal(t, acc, accById)
}

func TestDeleteAccount(t *testing.T) {
	defer cleanTestStore(t)

	acc := createRandomAccount(t, createRandomUser(t))
	accById, err := testStore.GetAccountById(ctx, acc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accById)

	acc, err = testStore.DeleteAccount(ctx, acc.ID)
	require.NoError(t, err)
	require.Equal(t, acc.ID, accById.ID)

	accById, err = testStore.GetAccountById(ctx, acc.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Empty(t, accById)
}

func TestGetTwoAccountsByIdForUpdate(t *testing.T) {
	defer cleanTestStore(t)

	var accs []Account
	var userIds []int64
	for i := 0; i < 2; i++ {
		usr := createRandomUser(t)
		userIds = append(userIds, usr.ID)
		accs = append(accs, createRandomAccount(t, usr))
	}

	gotAccs, err := testStore.GetTwoAccountsByIdForUpdate(ctx, GetTwoAccountsByIdForUpdateParams{userIds[0], userIds[1]})
	require.NoError(t, err)
	require.Equal(t, gotAccs[0], accs[0])
	require.Equal(t, gotAccs[1], accs[1])
}

func TestGetTwoAccountsByIdForUpdate_QueryError(t *testing.T) {
	params := GetTwoAccountsByIdForUpdateParams{}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, wantError)

	gotAccounts, gotError := New(dbMock).GetTwoAccountsByIdForUpdate(ctx, params)
	require.Nil(t, gotAccounts)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestGetTwoAccountsByIdForUpdate_ScanError(t *testing.T) {
	params := GetTwoAccountsByIdForUpdateParams{}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	rowsMock := new(rowsMock)
	rowsMock.On("Next").Return(true)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotAccounts, gotError := New(dbMock).GetTwoAccountsByIdForUpdate(ctx, params)
	require.Nil(t, gotAccounts)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Close")
}

func TestGetTwoAccountsByIdForUpdate_RowsError(t *testing.T) {
	params := GetTwoAccountsByIdForUpdateParams{}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	rowsMock := new(rowsMock)
	rowsMock.On("Next").Once().Return(true)
	rowsMock.On("Next").Once().Return(false)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	rowsMock.On("Err").Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotUsers, gotError := New(dbMock).GetTwoAccountsByIdForUpdate(ctx, params)
	require.Nil(t, gotUsers)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Err")
	rowsMock.AssertCalled(t, "Close")
}

func TestGetTotalCount(t *testing.T) {
	defer cleanTestStore(t)

	want := utils.GenRandIntInRange(5, 15)
	var wg sync.WaitGroup
	wg.Add(want)

	for i := 0; i < want; i++ {
		go func() {
			defer wg.Done()
			createRandomAccount(t, createRandomUser(t))
		}()
	}

	wg.Wait()
	got, err := testStore.GetTotalAccountsCount(ctx)
	require.NoError(t, err)
	require.Equal(t, got, int64(want))
}

func TestListAccounts(t *testing.T) {
	defer cleanTestStore(t)

	params := ListAccountsParams{Offset: 0, Limit: 2}
	got, err := testStore.ListAccounts(ctx, params)
	require.NoError(t, err)
	require.Empty(t, got)

	var want []ListAccountsRow
	for i := 0; i < 2; i++ {
		user := createRandomUser(t)
		acc := createRandomAccount(t, user)
		want = append(want, ListAccountsRow{
			ID:        acc.ID,
			UserID:    acc.UserID,
			FullName:  user.FullName,
			Email:     user.Email,
			Currency:  acc.Currency,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
		})
	}
	got, err = testStore.ListAccounts(ctx, params)
	require.NoError(t, err)
	require.Equal(t, got, want)
}

func TestListAccounts_QueryError(t *testing.T) {
	params := ListAccountsParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, wantError)

	gotAccounts, gotError := New(dbMock).ListAccounts(ctx, params)
	require.Nil(t, gotAccounts)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListAccounts_ScanError(t *testing.T) {
	params := ListAccountsParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	rowsMock := new(rowsMock)
	rowsMock.On("Next").Return(true)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotAccounts, gotError := New(dbMock).ListAccounts(ctx, params)
	require.Nil(t, gotAccounts)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Close")
}

func TestListAccounts_RowsError(t *testing.T) {
	params := ListAccountsParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	rowsMock := new(rowsMock)
	rowsMock.On("Next").Once().Return(true)
	rowsMock.On("Next").Once().Return(false)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	rowsMock.On("Err").Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotUsers, gotError := New(dbMock).ListAccounts(ctx, params)
	require.Nil(t, gotUsers)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Err")
	rowsMock.AssertCalled(t, "Close")
}

func TestOffsetBalance(t *testing.T) {
	defer cleanTestStore(t)

	acc := createRandomAccount(t, createRandomUser(t))
	deltaBalance := int64(utils.GenRandIntInRange(-1e6, 1e6))
	wantBalance := acc.Balance + deltaBalance

	gotAcc, err := testStore.OffsetAccountBalance(ctx, OffsetAccountBalanceParams{
		ID:    acc.ID,
		Delta: int64(deltaBalance),
	})

	require.NoError(t, err)
	require.Equal(t, wantBalance, gotAcc.Balance)
}
