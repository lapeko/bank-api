package repository

import (
	"context"
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T) *Entry {
	account := createTestAccount(t)
	params := CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.Int64(-1000, 1000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	return &entry
}

func createTestEntryByOwner(t *testing.T, ownerId int64) *Entry {
	params := CreateEntryParams{
		AccountID: ownerId,
		Amount:    random.Int64(-1000, 1000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	return &entry
}

func TestGetEntries(t *testing.T) {
	truncateTables()

	var entries []*Entry

	const size = 10

	for i := 0; i < size; i++ {
		entry := createTestEntry(t)
		entries = append(entries, entry)
	}

	limit := random.Int(1, size)
	offset := random.Int(0, size-limit)

	entries = entries[offset:]

	params := ListEntriesParams{Limit: int32(limit), Offset: int32(offset)}

	list, err := testQueries.ListEntries(context.Background(), params)

	require.NoError(t, err)
	require.NotNil(t, list)
	require.Equal(t, len(list), limit)

	for i, listItem := range list {
		require.NotNil(t, listItem)
		require.Equal(t, entries[i].ID, listItem.ID)
		require.Equal(t, entries[i].Amount, listItem.Amount)
		require.Equal(t, entries[i].AccountID, listItem.AccountID)
		require.Equal(t, entries[i].CreatedAt, listItem.CreatedAt)
	}
}

func TestGetEntriesByAccount(t *testing.T) {
	truncateTables()

	account := createTestAccount(t)
	otherAccount := createTestAccount(t)

	const size = 5
	for i := 0; i < size; i++ {
		createTestEntryByOwner(t, account.ID)
	}

	for i := 0; i < size; i++ {
		createTestEntryByOwner(t, otherAccount.ID)
	}

	params := ListEntriesByAccountParams{
		AccountID: account.ID,
		Offset:    0,
		Limit:     int32(size),
	}

	entries, err := testQueries.ListEntriesByAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotNil(t, entries)

	require.Equal(t, size, len(entries))

	for _, entry := range entries {
		require.Equal(t, account.ID, entry.AccountID)
		require.NotEmpty(t, entry.ID)
		require.NotEmpty(t, entry.CreatedAt)
	}
}
