package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostCategory struct {
	ID         int64 `db:"id"`
	PostID     int64 `db:"post_history_id"`
	CategoryID int64 `db:"category_id"`
	InsertTime int64 `db:"insert_time"`
}

func (database *Database) GetPostCategoryById(id int64) (*PostCategory, error) {
	cols := `id, post_id, category_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM category WHERE id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostCategoryById: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var pc PostCategory
	err = row.StructScan(&pc)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal category from GetPostCategoryById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &pc, nil
}

func (database *Database) GetPostCategories(post_history *PostHistory) ([]Category, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	c, err := database.getPostCategories(tx, post_history)
	if err != nil {
		fatal := "cannot get post categories in GetPostCategories: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetPostCategories: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostCategories: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return c, nil
}

func (database *Database) getPostCategories(tx *sqlx.Tx, post_history *PostHistory) ([]Category, error) {
	cols := `id, post_history_id, category_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post_categories WHERE post_history_id = $1`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostCategories: " + err.Error()
		rerr := tx.Rollback()
		if rerr != nil {
			fatal := "cannot rollback in getPostCategories: " + msg + ": " + rerr.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(post_history.ID)
	if err != nil {
		msg := "cannot execute query in getPostCategories: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in getPostCategories: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	pcs := []PostCategory{}
	for rows.Next() {
		var pc PostCategory
		err = rows.StructScan(&pc)
		if err != nil {
			msg := "cannot unmarshal category from getPostCategories: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in getPostCategories: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
		pcs = append(pcs, pc)
	}
	cs := []Category{}
	for i := range pcs {
		category, err := database.getCategoryByID(tx, pcs[i].CategoryID)
		if err != nil {
			msg := "cannot get category from getPostCategories: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in getPostCategories: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
		if category != nil {
			cs = append(cs, *category)
		}
	}
	return cs, nil
}

func (database *Database) insertPostCategories(tx *sqlx.Tx, post_history_id int64, category_ids []int64) ([]int64, error) {
	post_category_ids := []int64{}
	for i := range category_ids {
		category_id := category_ids[i]
		post_category_id, err := database.insertPostCategory(tx, post_history_id, category_id)
		if err != nil {
			msg := "cannot insert post category for insertPostCategories: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in insertPostCategories: " + msg + ": " + err.Error()
				return nil, errors.New(fatal)
			}
			return nil, errors.New(msg)
		}
		post_category_ids = append(post_category_ids, *post_category_id)

	}
	return post_category_ids, nil
}

func (database *Database) insertPostCategory(tx *sqlx.Tx, post_history_id int64, category_id int64) (*int64, error) {
	cols := `post_history_id, category_id`
	query := fmt.Sprintf(`INSERT INTO post_categories (%s) VALUES($1, $2)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPostCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post_history_id, category_id)
	if err != nil {
		msg := "cannot execute query in insertPostCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPostCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPostCategory but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	post_category_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPostCategory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostCategory: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return &post_category_id, nil
}
