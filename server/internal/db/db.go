package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteDriver = "sqlite3"
)

type DB struct {
	connStr string
}

func New(connStr string) *DB {
	return &DB{
		connStr: connStr,
	}
}

// Connect to the database using sqlite3
func (d *DB) Connect() *sql.DB {
	db, err := sql.Open(sqliteDriver, d.connStr)
	if err != nil {
		panic(err)
	}
	return db
}
