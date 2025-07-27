package db

import (
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils/testutils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	randAmount := int64(testutils.GenRandIntInRange(-1e6, 1e6))

	got, err := testStore.CreateEntry(ctx, CreateEntryParams{
		AccountID: acc.ID,
		Amount:    randAmount,
	})

	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.NotEmpty(t, got.ID)
	require.Equal(t, acc.ID, got.AccountID)
	require.Equal(t, got.Amount, randAmount)
	require.NotEmpty(t, got.CreatedAt)

	return got
}

func TestCreateEntry(t *testing.T) {
	defer cleanTestStore(t)

	createRandomEntry(t, createRandomAccount(t, createRandomUser(t)))
}

func TestGetEntryById(t *testing.T) {
	defer cleanTestStore(t)

	account := createRandomAccount(t, createRandomUser(t))
	entry := createRandomEntry(t, account)

	got, err := testStore.GetEntryById(ctx, entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, entry.ID, got.ID)
	require.Equal(t, entry.AccountID, got.AccountID)
	require.Equal(t, entry.Amount, got.Amount)
	require.Equal(t, entry.CreatedAt, got.CreatedAt)
}

func TestListEntries(t *testing.T) {
	defer cleanTestStore(t)

	account := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	entries, err := testStore.ListEntries(ctx, ListEntriesParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListEntriesByAccount(t *testing.T) {
	defer cleanTestStore(t)

	account := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 7; i++ {
		createRandomEntry(t, account)
	}

	entries, err := testStore.ListEntriesByAccount(ctx, ListEntriesByAccountParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	})
	require.NoError(t, err)
	require.True(t, len(entries) <= 5)
	for _, entry := range entries {
		require.Equal(t, account.ID, entry.AccountID)
	}
}

func TestGetTotalEntriesCount(t *testing.T) {
	defer cleanTestStore(t)

	account := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 3; i++ {
		createRandomEntry(t, account)
	}

	count, err := testStore.GetTotalEntriesCount(ctx)
	require.NoError(t, err)
	require.True(t, count >= 3)
}

func TestGetTotalEntriesCountByAccount(t *testing.T) {
	defer cleanTestStore(t)

	account := createRandomAccount(t, createRandomUser(t))
	for i := 0; i < 4; i++ {
		createRandomEntry(t, account)
	}

	count, err := testStore.GetTotalEntriesCountByAccount(ctx, account.ID)
	require.NoError(t, err)
	require.Equal(t, int64(4), count)
}
