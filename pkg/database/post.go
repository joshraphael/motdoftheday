package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/joshraphael/diary/pkg/post"
)

type Post struct {
	ID         int64  `db:"id"`
	UrlTitle   string `db:"url_title"`
	UserID     int64  `db:"user_id"`
	Title      string `db:"title"`
	Posted     int64  `db:"posted"`
	UpdateTime int64  `db:"update_time"`
	InsertTime int64  `db:"insert_time"`
}

type BOOL int64

const (
	db_TRUE  BOOL = 1
	db_FALSE BOOL = 0
)

func (database *Database) GetPostById(id int) (*Post, error) {
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
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

func (database *Database) GetPostByUrlTitle(url_title string) (*Post, error) {
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE LOWER(url_title) = LOWER($1)`, cols)
	stmt, err := database.db.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for GetPostByUrlTitle: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(url_title)
	var p Post
	err = row.StructScan(&p)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal post from GetPostByUrlTitle: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &p, nil
}

func (database *Database) SavePost(post post.Post) error {
	err := post.Validate()
	if err != nil {
		msg := "cannot validate post in SavePost: " + err.Error()
		return errors.New(msg)
	}
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for SavePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in SavePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	found, p, err := database.postExists(tx, post)
	if err != nil {
		msg := "cannot prepare statement for SavePost: " + err.Error()
		return errors.New(msg)
	}
	if BOOL(p.Posted) == db_TRUE {
		msg := "Post already posted and cannot be edited in SavePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in SavePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	var post_id int64
	if found {
		post_id = p.ID
		err = database.updatePost(tx, p, db_FALSE)
		if err != nil {
			msg := "cannot update post in SavePost: " + err.Error()
			return errors.New(msg)
		}
	} else {
		id, err := database.insertPost(tx, post, db_FALSE)
		if err != nil {
			msg := "cannot insert new post in SavePost: " + err.Error()
			return errors.New(msg)
		}
		post_id = id
	}
	_, err = database.insertPostHistory(tx, post_id, post)
	if err != nil {
		msg := "cannot insert post history in SavePost: " + err.Error()
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in SavePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in SavePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	return nil
}

func (database *Database) postExists(tx *sqlx.Tx, post post.Post) (bool, *Post, error) {
	url_title := strings.Join(strings.Split(post.Title, " "), "-")
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE LOWER(url_title) = LOWER($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for postExists: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in postExists: " + msg + ": " + err.Error()
			return false, nil, errors.New(fatal)
		}
		return false, nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(url_title)
	var p Post
	found := true
	err = row.StructScan(&p)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			found = false
		default:
			msg := "cannot unmarshal post from postExists: " + err.Error()
			return false, nil, errors.New(msg)
		}
	}
	return found, &p, nil
}

func (database *Database) insertPost(tx *sqlx.Tx, post post.Post, posted BOOL) (int64, error) {
	url_title := strings.Join(strings.Split(post.Title, " "), "-")
	cols := `url_title, user_id, title, posted`
	query := fmt.Sprintf(`INSERT INTO post (%s) VALUES($1, 1, $2, $3)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPost: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(url_title, post.Title, posted)
	if err != nil {
		msg := "cannot execute query in insertPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPost: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPost: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPost but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPost: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	post_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPost: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	return post_id, nil
}

func (database *Database) updatePost(tx *sqlx.Tx, db_post *Post, posted BOOL) error {
	query := `UPDATE post SET posted = $1, update_time = (CAST(strftime('%s', 'now') as integer)) WHERE id = $2`
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for updatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in updatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(posted, db_post.ID)
	if err != nil {
		msg := "cannot execute query in updatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in updatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in updatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in updatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in updatePost but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in updatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertPostHistory(tx *sqlx.Tx, post_id int64, post post.Post) (int64, error) {
	cols := `post_id, body`
	query := fmt.Sprintf(`INSERT INTO post_history (%s) VALUES($1, $2)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post_id, post.Body)
	if err != nil {
		msg := "cannot execute query in insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPostHistory but " + string(rows) + " rows were: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	post_history_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPostHistory: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in insertPostHistory: " + msg + ": " + err.Error()
			return -1, errors.New(fatal)
		}
		return -1, errors.New(msg)
	}
	return post_history_id, nil
}
