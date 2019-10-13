package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Tag struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	UserID     int    `db:"user_id"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetTagById(id int) (*Tag, error) {
	cols := `id, name, user_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetTagById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var t Tag
	err = row.StructScan(&t)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal tag from GetTagById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &t, nil
}

func (database *Database) GetTagByName(name string) (*Tag, error) {
	cols := `id, name, user_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM tag WHERE LOWER(name) = LOWER($1)`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetTagByName: " + err.Error()
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
			msg := "cannot unmarshal tag from GetTagByName: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &t, nil
}
