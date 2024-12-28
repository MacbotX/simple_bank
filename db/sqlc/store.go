package db

import (
	"context"
	"fmt"
	// "log"

	"github.com/jackc/pgx/v5"
)

// Store struct encapsulates the database and queries
type Store struct {
	db      *pgx.Conn
	*Queries
}

// NewStore creates a new Store instance
func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// exctx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
    // Start a transaction with default options
    tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return err
    }

    q := New(tx) // Initialize a new Queries object with the transaction
    err = fn(q)  // Execute the transactional function

    if err != nil {
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
        }
        return err
    }

    return tx.Commit(ctx) // Commit the transaction
}

type TransferTxParams struct {
	FromAccountID 	int64 	`json:"from_account_id"`
	ToAccountID 	int64 	`json:"to_account_id"`
	Amount 			int64 	`json:"amount"`
}

// TYransferTxResult is the result of the TransferTx transaction
type TransferTxResult struct {
	Transfer     	Transfer 	`json:"transfer"`
	FromAccount  	Account 	`json:"from_account"`
	ToAccount 		Account 	`json:"to_account"`
	FromEntry 		Entry 		`json:"from_entry"`
	ToEntry 		Entry 		`json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: 	arg.FromAccountID,
			ToAccountID: 	arg.ToAccountID,
			Amount: 		arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: 	arg.ToAccountID,
			Amount: 	arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update account balance

		return nil 
	})

	return result, err
}
