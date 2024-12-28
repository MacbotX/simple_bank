package db

import (
	"context"
	"testing"
	"time"

	"github.com/macbotxxx/simple_bank.git/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
			FromAccountID:    account1.ID,
			ToAccountID:    account2.ID,
			Amount:  util.RandomMoney(),
		}

	entry, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.FromAccountID, entry.FromAccountID)
	require.Equal(t, arg.ToAccountID, entry.ToAccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry

}


func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomTransfer(t)
	account2, err := testQueries.GetTransfer(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.FromAccountID, account2.FromAccountID)
	require.Equal(t, account1.ToAccountID, account2.ToAccountID)
	require.Equal(t, account1.Amount, account2.Amount)

	// Extract the time.Time values
	createdAt1 := account1.CreatedAt.Time
	createdAt2 := account2.CreatedAt.Time

	require.WithinDuration(t, createdAt1, createdAt2, time.Second)
}

func TestGetUpdateTransfer(t *testing.T) {
	account1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:      account1.ID,
		Amount: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Amount, account2.Amount)

	createdAt1 := account1.CreatedAt.Time
	createdAt2 := account2.CreatedAt.Time

	require.Equal(t, createdAt1, createdAt2, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	account1 := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetTransfer(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListTransfer(t *testing.T) {
	var lastAccount Transfer
	
	for i := 0; i < 10; i++ {
		lastAccount = createRandomTransfer(t)
	}

	arg := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, account := range transfers {
		require.NotEmpty(t, account)
		require.NotZero(t, account.ID)
		require.NotZero(t, account.CreatedAt)
	}

	require.LessOrEqual(t, transfers[0].ID, lastAccount.ID)
}
