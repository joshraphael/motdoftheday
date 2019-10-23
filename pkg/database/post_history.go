package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

type PostHistory struct {
	ID         int64  `db:"id"`
	PostID     int64  `db:"post_id"`
	Body       string `db:"body"`
	Method     string `db:"method"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetPostHistoryById(post_history_id int64) (*PostHistory, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetPostHistoryById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostHistoryById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	ph, err := database.getPostHistoryById(tx, post_history_id)
	if err != nil {
		msg := "cannot get post history in GetPostHistoryById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostHistoryById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetPostHistoryById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostHistoryById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return ph, nil
}

func (database *Database) GetLatestPostHistory(post *Post) (*PostHistory, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetLatestPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetLatestPost: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	ph, err := database.getLatestPost(tx, post)
	if err != nil {
		msg := "cannot get post history in GetLatestPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetLatestPost: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetLatestPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetLatestPost: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return ph, nil
}

func (database *Database) getLatestPost(tx *sqlx.Tx, post *Post) (*PostHistory, error) {
	cols := `MAX(id) AS id, post_id, body, method, insert_time`
	query := fmt.Sprintf(`
	SELECT %s
	FROM post_history
	WHERE insert_time = (
	  SELECT MAX(insert_time)
	  FROM post_history
	  WHERE post_id = $1
	)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getLatestPost: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(post.ID)
	var ph PostHistory
	err = row.StructScan(&ph)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := "no post history found for this post: " + err.Error()
			return nil, errors.New(msg)
		default:
			msg := "cannot unmarshal post from getLatestPost: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &ph, nil
}

func (database *Database) getPostHistory(tx *sqlx.Tx, post *Post) ([]PostHistory, error) {
	cols := `id, post_id, body, method, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post_history WHERE post_id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getLatestPost: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(post.ID)
	if err != nil {
		msg := "cannot query statement for getLatestPost: " + err.Error()
		return nil, errors.New(msg)
	}
	phs := []PostHistory{}
	for rows.Next() {
		var ph PostHistory
		err = rows.StructScan(&ph)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				msg := "no post history found for this post: " + err.Error()
				return nil, errors.New(msg)
			default:
				msg := "cannot unmarshal post from getLatestPost: " + err.Error()
				return nil, errors.New(msg)
			}
		}
		phs = append(phs, ph)
	}
	return phs, nil
}

func (database *Database) insertPostHistory(tx *sqlx.Tx, post_id int64, post post.Post) (*int64, error) {
	cols := `post_id, body, method`
	query := fmt.Sprintf(`INSERT INTO post_history (%s) VALUES($1, $2, $3)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPostHistory: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post_id, post.Body, post.Method())
	if err != nil {
		msg := "cannot execute query in insertPostHistory: " + err.Error()
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPostHistory: " + err.Error()
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPostHistory but " + string(rows) + " rows were: " + err.Error()
		return nil, errors.New(msg)
	}
	post_history_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPostHistory: " + err.Error()
		return nil, errors.New(msg)
	}
	return &post_history_id, nil
}

func (database *Database) getPostHistoryById(tx *sqlx.Tx, post_history_id int64) (*PostHistory, error) {
	cols := `id, post_id, body, method, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post_history WHERE id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostHistoryById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(post_history_id)
	var ph PostHistory
	err = row.StructScan(&ph)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := "no post history found for this post: " + err.Error()
			return nil, errors.New(msg)
		default:
			msg := "cannot unmarshal post from getLatestPost: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &ph, nil
}
