package database

import (
	"database/sql"
	"errors"
)

type Database struct {
	db *sql.DB
}

func New(database *sql.DB) (*Database, error) {
	err := database.Ping()
	if err != nil {
		msg := "bad ping: " + err.Error()
		return nil, errors.New(msg)
	}
	return &Database{
		db: database,
	}, nil
}
