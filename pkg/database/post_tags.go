package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type PostTag struct {
	ID         int   `db:"id"`
	PostID     int   `db:"post_id"`
	TagID      int   `db:"tag_id"`
	InsertTime int64 `db:"insert_time"`
}

func (database *Database) GetPostTagById(id int) (*PostTag, error) {
	cols := `id, post_id, tag_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostTagById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var pt PostTag
	err = row.StructScan(&pt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal tag from GetPostTagById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &pt, nil
}

func (database *Database) GetPostTagByPostId(postId int) (*PostTag, error) {
	cols := `id, post_id, tag_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE post_id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostTagById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(postId)
	var pt PostTag
	err = row.StructScan(&pt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal tag from GetPostTagById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &pt, nil
}
