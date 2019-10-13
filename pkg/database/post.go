package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Post struct {
	ID         int    `db:"id"`
	UserID     int    `db:"user_id"`
	Title      string `db:"title"`
	Body       string `db:"body"`
	PostTime   int64  `db:"post_time"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetPostById(id int) (*Post, error) {
	cols := `id, user_id, title, body, post_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var p Post
	err = row.StructScan(&p)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal post from GetPostById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &p, nil
}

func (database *Database) GetPostByTitle(title string) (*Post, error) {
	cols := `id, user_id, title, body, post_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE LOWER(title) = LOWER($1)`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostByTitle: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(title)
	var p Post
	err = row.StructScan(&p)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal post from GetPostByTitle: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &p, nil
}
