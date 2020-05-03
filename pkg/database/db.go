package database

import (
	"errors"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db  *sqlx.DB
	cfg Config
}

type CompletePost struct {
	Post       *Post
	History    []PostHistory
	Categories map[int64][]Category
	Tags       map[int64][]Tag
}

type CompletePostHistory struct {
	Post       *Post
	History    *PostHistory
	Categories []Category
	Tags       []Tag
}

func New(c Config) (*Database, *sqlx.DB, error) {
	db_name := "./" + c.File
	if _, err := os.Stat(db_name); err != nil {
		msg := "Database " + db_name + " does not exist: " + err.Error()
		log.Fatalln(msg)
	}
	database, err := sqlx.Open("sqlite3", db_name+"?_foreign_keys=on")
	if err != nil {
		log.Fatalln(err)
	}
	err = database.Ping()
	if err != nil {
		msg := "bad ping: " + err.Error()
		return nil, nil, errors.New(msg)
	}
	return &Database{
		db:  database,
		cfg: c,
	}, database, nil
}
