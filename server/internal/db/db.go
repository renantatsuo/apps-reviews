package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteDriver = "sqlite3"
	sqliteFile   = "file:data/database.db?cache=shared&mode=rwc"
)

// Connect to the database using sqlite3
func Connect() *sql.DB {
	db, err := sql.Open(sqliteDriver, sqliteFile)
	if err != nil {
		panic(err)
	}
	return db
}
