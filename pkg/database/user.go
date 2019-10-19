package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID         int64  `db:"id"`
	Username   string `db:"user_name"`
	Firstname  string `db:"first_name"`
	Lastname   string `db:"last_name"`
	UpdateTime int64  `db:"update_time"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetUserById(id int64) (*User, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetUserById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetUserById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	u, err := database.getUserById(tx, id)
	if err != nil {
		msg := "cannot get user in GetUserById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetUserById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetUserById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetUserById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getUserById(tx *sqlx.Tx, id int64) (*User, error) {
	cols := `id, user_name, first_name, last_name, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM user WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getUserById: " + err.Error()
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
			msg := "cannot unmarshal user from getUserById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &u, nil
}

func (database *Database) GetUserByUsername(username string) (*User, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetUserByUsername: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetUserByUsername: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	u, err := database.getUserByUsername(tx, username)
	if err != nil {
		msg := "cannot get user in GetUserByUsername: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetUserByUsername: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetUserByUsername: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetUserByUsername: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getUserByUsername(tx *sqlx.Tx, username string) (*User, error) {
	cols := `id, user_name, first_name, last_name, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM user WHERE LOWER(user_name) = LOWER($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getUserByUsername: " + err.Error()
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
			msg := "cannot unmarshal user from getUserByUsername: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &u, nil
}
