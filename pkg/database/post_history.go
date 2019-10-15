package database

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

func (database *Database) insertPostHistory(tx *sqlx.Tx, post_id int64, post post.Post) (*int64, error) {
	cols := `post_id, body, method`
	query := fmt.Sprintf(`INSERT INTO post_history (%s) VALUES($1, $2, $3)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post_id, post.Body, post.Method())
	if err != nil {
		msg := "cannot execute query in insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPostHistory but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	post_history_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return &post_history_id, nil
}
