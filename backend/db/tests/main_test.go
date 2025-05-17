package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	db "github/kasho/backend/db/sqlc"
)

var testQuery *db.Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgresql://root:root@localhost:5432/kasho_db?sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	testQuery = db.New(conn)

	os.Exit(m.Run())
}