package db

import (
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/db/utils"
	"github.com/stretchr/testify/require"
)

func TestTransferMoney(t *testing.T) {
	defer cleanTestStore(t)

	u1 := createRandomUser(t)
	u2 := createRandomUser(t)
	acc1Balance := utils.RandInt64InRange(1e4, 1e5)
	acc2Balance := utils.RandInt64InRange(1e4, 1e5)
	transferAmmount := utils.RandInt64InRange(1, 1e4)
	acc1 := createAccountWithParams(t, CreateAccountParams{UserID: u1.ID, Currency: utils.CurrencyUSD, Balance: acc1Balance})
	acc2 := createAccountWithParams(t, CreateAccountParams{UserID: u2.ID, Currency: utils.CurrencyUSD, Balance: acc2Balance})

	err := testStore.TransferMoney(ctx, acc1.ID, acc2.ID, transferAmmount)
	require.NoError(t, err)

	acc1, err = testStore.GetQueries().GetAccountById(ctx, acc1.ID)
	require.NoError(t, err)
	acc2, err = testStore.GetQueries().GetAccountById(ctx, acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, acc1Balance-transferAmmount)
	require.Equal(t, acc2.Balance, acc2Balance+transferAmmount)
}
