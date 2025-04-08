package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	firstAccount := createRandomAccount(t)
	secondAccount := createRandomAccount(t)

	concurrentTransfers := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < concurrentTransfers; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: firstAccount.ID,
				ToAccountID:   secondAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	// check results
	for i := 0; i < concurrentTransfers; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, firstAccount.ID, transfer.FromAccountID)
		require.Equal(t, secondAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, firstAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, secondAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, firstAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, secondAccount.ID, toAccount.ID)

		firstDifference := firstAccount.Balance - fromAccount.Balance
		secondDifference := toAccount.Balance - secondAccount.Balance

		require.Equal(t, firstDifference, secondDifference)
		require.True(t, firstDifference > 0)
		require.True(t, firstDifference%amount == 0) // amount must be multiple of concurrentTransfers

		k := int(firstDifference / amount)
		require.True(t, k >= 1 && k <= concurrentTransfers)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedFirstAccount, err := testQueries.GetAccount(context.Background(), firstAccount.ID)
	require.NoError(t, err)

	updatedSecondAccount, err := testQueries.GetAccount(context.Background(), secondAccount.ID)
	require.NoError(t, err)

	require.Equal(t, firstAccount.Balance-int64(concurrentTransfers)*amount, updatedFirstAccount.Balance)
	require.Equal(t, secondAccount.Balance+int64(concurrentTransfers)*amount, updatedSecondAccount.Balance)
}
