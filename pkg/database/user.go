package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID         int    `db:"id"`
	Username   string `db:"user_name"`
	Firstname  string `db:"first_name"`
	Lastname   string `db:"last_name"`
	UpdateTime int64  `db:"update_time"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetUserById(id int) (*User, error) {
	cols := `id, user_name, first_name, last_name, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM user WHERE id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetUserById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var u User
	err = row.StructScan(&u)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from GetUserById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &u, nil
}

func (database *Database) GetUserByUsername(username string) (*User, error) {
	cols := `id, user_name, first_name, last_name, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM user WHERE LOWER(user_name) = LOWER($1)`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetUserByUsername: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(username)
	var u User
	err = row.StructScan(&u)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from GetUserByUsername: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &u, nil
}
