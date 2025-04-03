package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store struct encapsulates the database and queries
type Store struct {
	db      *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store instance
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db), // Assuming New initializes queries from a pool or transaction
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
    // Start a transaction with default options
    tx, err := store.db.BeginTx(ctx, pgx.TxOptions{}) // Ensure correct options usage from pgx/v5
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    // Initialize a new Queries object with the transaction
    q := New(tx) // Initialize with the transaction, not the pool
    err = fn(q)  // Execute the transactional function

    if err != nil {
        // Attempt to rollback the transaction on error
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
        }
        return err
    }

    // Commit the transaction on success
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
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
		

		return nil 
	})

	return result, err
}
