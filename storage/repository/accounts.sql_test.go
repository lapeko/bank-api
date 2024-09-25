package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) *Account {
	params := CreateAccountParams{
		Owner:    random.String(6),
		Currency: random.Currency(),
		Balance:  random.Int64(0, 1000),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)
	require.NotEmpty(t, account.ID)
	require.NotEmpty(t, account.CreatedAt)

	return &account
}

func TestGetAccounts(t *testing.T) {
	var accounts []*Account

	const size = 10

	for i := 0; i < size; i++ {
		account := createTestAccount(t)
		accounts = append(accounts, account)
	}

	limit := random.Int(1, size)
	offset := random.Int(0, int(size-limit))

	params := ListAccountsParams{Limit: int32(limit), Offset: int32(offset)}

	list, err := testQueries.ListAccounts(context.Background(), params)

	require.NoError(t, err)
	require.NotNil(t, list)
	require.Equal(t, len(list), limit)

	for i, listItem := range list {
		require.NotNil(t, listItem)
		require.Equal(t, listItem.ID, accounts[offset+i].ID)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createTestAccount(t)
	params := UpdateAccountParams{ID: account.ID, Balance: random.Int64(0, 1000000)}
	err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)
	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedAccount)
	require.Equal(t, params.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createTestAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	emptyAcc, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Empty(t, emptyAcc)
	require.NotEmpty(t, err)
	require.Equal(t, err, sql.ErrNoRows)
}
