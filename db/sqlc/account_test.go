package db

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/db/utils"
	"github.com/stretchr/testify/require"
)

func randCurrency() utils.Currency {
	switch rand.Int() % 3 {
	case 0:
		return utils.CurrencyEURO
	case 1:
		return utils.CurrencyPLN
	default:
		return utils.CurrencyUSD
	}
}

func init() {

}

func TestAccount(t *testing.T) {
	user := createRandomUser(t)

	want := CreateAccountParams{
		UserID:   user.ID,
		Currency: randCurrency(),
		Balance:  rand.Int63(),
	}

	fmt.Println()
	fmt.Println(want.Currency)

	got, err := sqlcQueries.CreateAccount(context.TODO(), want)

	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.NotZero(t, got.ID)
	require.Equal(t, got.UserID, want.UserID)
	require.Equal(t, got.Currency, want.Currency)
	require.Equal(t, got.Balance, want.Balance)
	require.NotZero(t, got.CreatedAt)
}
