package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

type Category struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	UserID     int64  `db:"user_id"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetCategoryById(category_id int64) (*Category, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetCategoryById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCategoryById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	c, err := database.getCategoryByID(tx, category_id)
	if err != nil {
		fatal := "cannot get category in GetCategoryById: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetCategoryById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCategoryById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return c, nil
}

func (database *Database) GetCategoryByName(name string) (*Category, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetCategoryByName: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCategoryByName: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	c, err := database.getCategoryByName(tx, name)
	if err != nil {
		fatal := "cannot get category in GetCategoryByName: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetCategoryByName: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCategoryByName: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return c, nil
}

func (database *Database) getCategoryByName(tx *sqlx.Tx, name string) (*Category, error) {
	cols := `id, name, user_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM category WHERE LOWER(name) = LOWER($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getCategoryByName: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in getCategoryByName: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(name)
	var c Category
	err = row.StructScan(&c)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal category from getCategoryByName: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in getCategoryByName: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
	}
	return &c, nil
}

func (database *Database) getCategoryByID(tx *sqlx.Tx, category_id int64) (*Category, error) {
	cols := `id, name, user_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM category WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getCategoryByID: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in getCategoryByID: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(category_id)
	var c Category
	err = row.StructScan(&c)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal category from getCategoryByID: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in getCategoryByID: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
	}
	return &c, nil
}

func (database *Database) insertCategories(tx *sqlx.Tx, post post.Post) ([]int64, error) {
	url_categories := post.UrlCategories()
	category_ids := []int64{}
	for i := range url_categories {
		name := url_categories[i]
		category, err := database.getCategoryByName(tx, name)
		if err != nil {
			msg := "cannot get category for insertCategories: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in insertCategories: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
		if category != nil {
			category_ids = append(category_ids, category.ID)
		} else {
			category_id, err := database.insertCategory(tx, name)
			if err != nil {
				msg := "cannot insert category for insertCategories: " + err.Error()
				err = tx.Rollback()
				if err != nil {
					fatal := "cannot rollback in insertCategories: " + msg + ": " + err.Error()
					return nil, errors.New(fatal)
				}
				return nil, errors.New(msg)
			}
			category_ids = append(category_ids, *category_id)
		}
	}
	return category_ids, nil
}

func (database *Database) insertCategory(tx *sqlx.Tx, name string) (*int64, error) {
	cols := `user_id, name`
	query := fmt.Sprintf(`INSERT INTO category (%s) VALUES(1, $1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(name)
	if err != nil {
		msg := "cannot execute query in insertCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertCategory but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	category_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return &category_id, nil
}
