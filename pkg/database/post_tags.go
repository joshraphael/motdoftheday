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

func (database *Database) GetPostTagById(id int64) (*PostTag, error) {
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

func (database *Database) GetPostTagsByPostId(post_id int64) ([]PostTag, error) {
	cols := `id, post_id, tag_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE post_id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostTagById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(post_id)
	if err != nil {

	}
	pts := []PostTag{}
	for rows.Next() {
		var pt PostTag
		err = rows.StructScan(&pt)
		if err != nil {
			msg := "cannot unmarshal tag from GetPostTagById: " + err.Error()
			return nil, errors.New(msg)
		}
		pts = append(pts, pt)
	}
	return pts, nil
}

func (database *Database) insertPostTags(tx *sqlx.Tx, post_history_id int64, tag_ids []int64) ([]int64, error) {
	post_tag_ids := []int64{}
	for i := range tag_ids {
		tag_id := tag_ids[i]
		post_tag_id, err := database.insertPostTag(tx, post_history_id, tag_id)
		if err != nil {
			msg := "cannot insert post tag for insertPostTags: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in insertPostTags: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
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
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post_history_id, tag_id)
	if err != nil {
		msg := "cannot execute query in insertPostTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPostTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPostTag but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	post_tag_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPostTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return &post_tag_id, nil
}
