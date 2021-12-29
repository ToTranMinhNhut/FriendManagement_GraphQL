package config

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var (
	dbUrl        = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	maxOpenConns = 10
	maxIdleConns = 5
)

func NewDatabase() (*sql.DB, error) {
	if dbUrl == "" {
		return nil, errors.New("DATABASE_URL not found")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}

func CloseDatabase(db *sql.DB) error {
	return db.Close()
}
