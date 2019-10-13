package database

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func New(database *sqlx.DB) (*Database, error) {
	err := database.Ping()
	if err != nil {
		msg := "bad ping: " + err.Error()
		return nil, errors.New(msg)
	}
	return &Database{
		db: database,
	}, nil
}
