package db

import (
	"context"
	"log"
	"os"
	"testing"

	// "github.com/jackc/pgx/v5" for single connection
	"github.com/jackc/pgx/v5/pgxpool" // for connection pool
)


// pgxpool.Pool is a connection pool, not a single connection (*pgx.Conn). 
// Use it when you need a pool of connections shared across your application or tests.
// If you're assigning pgxpool.Pool to a variable expecting *pgx.Conn, the will be error.


const (
	dbSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDB  *pgxpool.Pool
)


func TestMain(m *testing.M) {
	var err error

	// Create a connection pool for tests
	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer testDB.Close()

	// Initialize queries with the database connection
	testQueries = New(testDB)

	// Run tests
	os.Exit(m.Run())
}
