package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

type Tag struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	UserID     int64  `db:"user_id"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetTagById(tag_id int64) (*Tag, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetTagById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetTagById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	t, err := database.getTagByID(tx, tag_id)
	if err != nil {
		fatal := "cannot get tag in GetTagById: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetTagById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetTagById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return t, nil
}

func (database *Database) GetTagByName(name string) (*Tag, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetTagByName: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetTagByName: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	t, err := database.getTagByName(tx, name)
	if err != nil {
		fatal := "cannot get tag in GetTagByName: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetTagByName: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetTagByName: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return t, nil
}

func (database *Database) getTagByName(tx *sqlx.Tx, name string) (*Tag, error) {
	cols := `id, name, user_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE LOWER(name) = LOWER($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getTagByName: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in getTagByName: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(name)
	var t Tag
	err = row.StructScan(&t)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal tag from getTagByName: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in getTagByName: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
	}
	return &t, nil
}

func (database *Database) getTagByID(tx *sqlx.Tx, tag_id int64) (*Tag, error) {
	cols := `id, name, user_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getTagByID: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in getTagByID: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(tag_id)
	var t Tag
	err = row.StructScan(&t)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal tag from getTagByID: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in getTagByID: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
	}
	return &t, nil
}

func (database *Database) insertTags(tx *sqlx.Tx, post post.Post) ([]int64, error) {
	url_tags := post.UrlTags()
	tag_ids := []int64{}
	for i := range url_tags {
		name := url_tags[i]
		tag, err := database.getTagByName(tx, name)
		if err != nil {
			msg := "cannot get tag for insertTags: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in insertTags: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
		if tag != nil {
			tag_ids = append(tag_ids, tag.ID)
		} else {
			tag_id, err := database.insertTag(tx, name)
			if err != nil {
				msg := "cannot insert tag for insertTags: " + err.Error()
				err = tx.Rollback()
				if err != nil {
					fatal := "cannot rollback in insertTags: " + msg + ": " + err.Error()
					return nil, errors.New(fatal)
				}
				return nil, errors.New(msg)
			}
			tag_ids = append(tag_ids, *tag_id)
		}
	}
	return tag_ids, nil
}

func (database *Database) insertTag(tx *sqlx.Tx, name string) (*int64, error) {
	cols := `user_id, name`
	query := fmt.Sprintf(`INSERT INTO tag (%s) VALUES(1, $1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(name)
	if err != nil {
		msg := "cannot execute query in insertTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertTag but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	tag_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertTag: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertTag: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return &tag_id, nil
}
