package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/stretchr/testify/require"
)

func deleteAccounts(t *testing.T) {
	_, err := testDb.Exec("DELETE FROM accounts")
	require.NoError(t, err)
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{Offset: 0, Limit: 1})
	require.Empty(t, err)
	require.Empty(t, accounts)
}

func createTestAccount(t *testing.T) *Account {
	params := CreateAccountParams{
		Owner:    random.String(6),
		Currency: random.Currency(),
		Balance:  random.Int64(0, 1000),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, params.Owner)
	require.Equal(t, account.Balance, params.Balance)
	require.Equal(t, account.Currency, params.Currency)
	require.NotEmpty(t, account.ID)
	require.NotEmpty(t, account.CreatedAt)

	return &account
}

func createTestAccountWithBalance(t *testing.T, balance int64) *Account {
	params := CreateAccountParams{
		Owner:    random.String(6),
		Currency: random.Currency(),
		Balance:  balance,
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, params.Owner)
	require.Equal(t, account.Balance, params.Balance)
	require.Equal(t, account.Currency, params.Currency)
	require.NotEmpty(t, account.ID)
	require.NotEmpty(t, account.CreatedAt)

	return &account
}

func TestGetAccounts(t *testing.T) {
	deleteAccounts(t)

	var accounts []*Account

	const size = 10

	for i := 0; i < size; i++ {
		account := createTestAccount(t)
		accounts = append(accounts, account)
	}

	limit := random.Int(1, size)
	offset := random.Int(0, size-limit)

	accounts = accounts[offset:]

	params := ListAccountsParams{Limit: int32(limit), Offset: int32(offset)}

	list, err := testQueries.ListAccounts(context.Background(), params)

	require.NoError(t, err)
	require.NotNil(t, list)
	require.Equal(t, len(list), limit)

	for i, listItem := range list {
		require.NotNil(t, listItem)
		require.Equal(t, accounts[i].ID, listItem.ID)
		require.Equal(t, accounts[i].Balance, listItem.Balance)
		require.Equal(t, accounts[i].Currency, listItem.Currency)
		require.Equal(t, accounts[i].Owner, listItem.Owner)
		require.Equal(t, accounts[i].CreatedAt, listItem.CreatedAt)
	}
}

func TestUpdateAccount(t *testing.T) {
	deleteAccounts(t)
	account := createTestAccount(t)
	params := UpdateAccountParams{ID: account.ID, Balance: random.Int64(0, 1000000)}
	updatedAccount, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NoError(t, err)
	require.NotNil(t, updatedAccount)
	require.Equal(t, params.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	deleteAccounts(t)
	account := createTestAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	emptyAcc, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Empty(t, emptyAcc)
	require.NotEmpty(t, err)
	require.Equal(t, err, sql.ErrNoRows)
}
