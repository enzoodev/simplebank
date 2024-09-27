package db

import (
	"context"
	"testing"
	"time"

	"github.com/enzoodev/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	
	arg := CreateEntryParams{
		AccountID: util.Int64ToNullInt64(account.ID),
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			AccountID: util.Int64ToNullInt64(account.ID),
			Amount: util.RandomMoney(),
		}

		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListEntriesParams{
		AccountID: util.Int64ToNullInt64(account.ID),
		Limit: 5,
		Offset: 0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.True(t, entry.AccountID.Int64 == account.ID)
	}
}