package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostTag struct {
	ID         int64 `db:"id"`
	PostID     int64 `db:"post_history_id"`
	TagID      int64 `db:"tag_id"`
	InsertTime int64 `db:"insert_time"`
}

func (database *Database) getPostTagById(tx *sqlx.Tx, id int64) (*PostTag, error) {
	cols := `id, post_id, tag_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostTagById: " + err.Error()
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
			msg := "cannot unmarshal tag from getPostTagById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &pt, nil
}

func (database *Database) GetPostTags(post_history *PostHistory) ([]Tag, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for GetPostTags: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostTags: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	t, err := database.getPostTags(tx, post_history)
	if err != nil {
		msg := "cannot get post tags in GetPostTags: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostTags: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetPostTags: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostTags: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return t, nil
}

func (database *Database) getPostTags(tx *sqlx.Tx, post_history *PostHistory) ([]Tag, error) {
	cols := `id, post_history_id, tag_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post_tags WHERE post_history_id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostTags: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(post_history.ID)
	if err != nil {
		msg := "cannot execute query in getPostTags: " + err.Error()
		return nil, errors.New(msg)
	}
	pts := []PostTag{}
	for rows.Next() {
		var pt PostTag
		err = rows.StructScan(&pt)
		if err != nil {
			msg := "cannot unmarshal tag from getPostTags: " + err.Error()
			return nil, errors.New(msg)
		}
		pts = append(pts, pt)
	}
	ts := []Tag{}
	for i := range pts {
		tag, err := database.getTagByID(tx, pts[i].TagID)
		if err != nil {
			msg := "cannot get tag from getPostTags: " + err.Error()
			return nil, errors.New(msg)
		}
		if tag != nil {
			ts = append(ts, *tag)
		}
	}
	return ts, nil
}

func (database *Database) insertPostTags(tx *sqlx.Tx, post_history_id int64, tag_ids []int64) ([]int64, error) {
	post_tag_ids := []int64{}
	for i := range tag_ids {
		tag_id := tag_ids[i]
		post_tag_id, err := database.insertPostTag(tx, post_history_id, tag_id)
		if err != nil {
			msg := "cannot insert post tag for insertPostTags: " + err.Error()
			return nil, errors.New(msg)
		}
		post_tag_ids = append(post_tag_ids, *post_tag_id)

	}
	return post_tag_ids, nil
}

func (database *Database) insertPostTag(tx *sqlx.Tx, post_history_id int64, tag_id int64) (*int64, error) {
	cols := `post_history_id, tag_id`
	query := fmt.Sprintf(`INSERT INTO post_tags (%s) VALUES($1, $2)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPostTag: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post_history_id, tag_id)
	if err != nil {
		msg := "cannot execute query in insertPostTag: " + err.Error()
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPostTag: " + err.Error()
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPostTag but " + string(rows) + " rows were: " + err.Error()
		return nil, errors.New(msg)
	}
	post_tag_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPostTag: " + err.Error()
		return nil, errors.New(msg)
	}
	return &post_tag_id, nil
}
