package db

import (
	"context"
	"testing"
	"time"

	"github.com/macbotxxx/simple_bank.git/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccount(t)

	arg := CreateEntryParams{
			AccountID:    account1.ID,
			Amount:  util.RandomMoney(),
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
	account1 := createRandomEntry(t)
	account2, err := testQueries.GetEntry(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, account1.Amount, account2.Amount)

	// Extract the time.Time values
	createdAt1 := account1.CreatedAt.Time
	createdAt2 := account2.CreatedAt.Time

	require.WithinDuration(t, createdAt1, createdAt2, time.Second)
}

func TestGetUpdateEntry(t *testing.T) {
	account1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:      account1.ID,
		Amount: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, arg.Amount, account2.Amount)


	createdAt1 := account1.CreatedAt.Time
	createdAt2 := account2.CreatedAt.Time

	require.Equal(t, createdAt1, createdAt2, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	account1 := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetEntry(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListEntry(t *testing.T) {
	var lastAccount Entry
	
	for i := 0; i < 10; i++ {
		lastAccount = createRandomEntry(t)
	}

	arg := ListEntryParams{
		Limit:  5,
		Offset: 5,
	}

	entires, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entires, 5)

	for _, account := range entires {
		require.NotEmpty(t, account)
		require.NotZero(t, account.ID)
		require.NotZero(t, account.CreatedAt)
	}

	require.LessOrEqual(t, entires[0].ID, lastAccount.ID)
}
