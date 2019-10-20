package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
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

func (b BOOL) Value() int64 {
	return int64(b)
}

const (
	db_TRUE  BOOL = 1
	db_FALSE BOOL = 0
)

func DB_TRUE() BOOL {
	return db_TRUE
}

func DB_FALSE() BOOL {
	return db_FALSE
}

func (database *Database) GetPostById(id int64) (*Post, error) {
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
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for GetPostByUrlTitle: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostByUrlTitle: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	p, err := database.getPostByUrlTitle(tx, url_title)
	if err != nil {
		fatal := "cannot get post in GetPostByUrlTitle: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetPostByUrlTitle: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetPostByUrlTitle: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return p, nil
}

func (database *Database) GetDraftPosts() ([]Post, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for GetDraftPosts: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetDraftPosts: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	ps, err := database.getPostsByPosted(tx, DB_FALSE())
	if err != nil {
		fatal := "cannot get drafts in GetDraftPosts: " + err.Error()
		return nil, errors.New(fatal)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetDraftPosts: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetDraftPosts: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	return ps, nil
}

func (database *Database) GetCompletePost(post *Post) (*CompletePost, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for GetCompletePostById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCompletePostById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	history, err := database.getPostHistory(tx, post)
	if err != nil {
		msg := "cannot get post history for GetCompletePostById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCompletePostById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	post_categories, err := database.getPostCategoriesByHistory(tx, history)
	if err != nil {
		msg := "cannot get post history categories for GetCompletePostById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCompletePostById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	post_tags, err := database.getPostTagsByHistory(tx, history)
	if err != nil {
		msg := "cannot get post history tags for GetCompletePostById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCompletePostById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetCompletePostById: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in GetCompletePostById: " + msg + ": " + err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	complete_post := &CompletePost{
		Post:       post,
		History:    history,
		Categories: post_categories,
		Tags:       post_tags,
	}
	return complete_post, nil
}

func (database *Database) CreatePost(post post.Post, posted BOOL) error {
	err := post.Validate()
	if err != nil {
		msg := "cannot validate post in CreatePost: " + err.Error()
		return errors.New(msg)
	}
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "begin transaction for CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	found, p, err := database.getPost(tx, post)
	if err != nil {
		msg := "cannot prepare statement for CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	if BOOL(p.Posted) == db_TRUE {
		msg := "Post already posted and cannot be edited in CreatePost"
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	var post_id int64
	if found {
		err = database.updatePost(tx, p, posted)
		if err != nil {
			msg := "cannot update post in CreatePost: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
				return errors.New(fatal)
			}
			return errors.New(msg)
		}
		post_id = p.ID
	} else {
		id, err := database.insertPost(tx, post, posted)
		if err != nil {
			msg := "cannot insert new post in CreatePost: " + err.Error()
			err = tx.Rollback()
			if err != nil {
				fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
				return errors.New(fatal)
			}
			return errors.New(msg)
		}
		post_id = *id
	}
	post_history_id, err := database.insertPostHistory(tx, post_id, post)
	if err != nil {
		msg := "cannot insert post history in CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	category_ids, err := database.insertCategories(tx, post)
	if err != nil {
		msg := "cannot insert categories in CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	tag_ids, err := database.insertTags(tx, post)
	if err != nil {
		msg := "cannot insert tags in CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	_, err = database.insertPostCategories(tx, *post_history_id, category_ids)
	if err != nil {
		msg := "cannot insert post categories in CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	_, err = database.insertPostTags(tx, *post_history_id, tag_ids)
	if err != nil {
		msg := "cannot insert post tags in CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in CreatePost: " + err.Error()
		err = tx.Rollback()
		if err != nil {
			fatal := "cannot rollback in CreatePost: " + msg + ": " + err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	return nil
}

func (database *Database) getPostByUrlTitle(tx *sqlx.Tx, url_title string) (*Post, error) {
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE LOWER(url_title) = LOWER($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostByUrlTitle: " + err.Error()
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
			msg := "cannot unmarshal post from getPostByUrlTitle: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &p, nil
}

func (database *Database) getPostById(tx *sqlx.Tx, id int64) (*Post, error) {
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostById: " + err.Error()
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
			msg := "cannot unmarshal post from getPostById: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &p, nil
}

func (database *Database) getPost(tx *sqlx.Tx, post post.Post) (bool, *Post, error) {
	url_title := post.UrlTitle()
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE LOWER(url_title) = LOWER($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for postExists: " + err.Error()
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

func (database *Database) getPostsByPosted(tx *sqlx.Tx, posted BOOL) ([]Post, error) {
	cols := `id, url_title, user_id, title, posted, update_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM post WHERE posted = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for getPostsByPosted: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(posted)
	if err != nil {
		msg := "cannot execute statement for getPostsByPosted: " + err.Error()
		return nil, errors.New(msg)
	}
	ps := []Post{}
	for rows.Next() {
		var p Post
		err = rows.StructScan(&p)
		if err != nil {
			msg := "cannot unmarshal post from getPostsByPosted: " + err.Error()
			return nil, errors.New(msg)
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (database *Database) insertPost(tx *sqlx.Tx, post post.Post, posted BOOL) (*int64, error) {
	url_title := post.UrlTitle()
	cols := `url_title, user_id, title, posted`
	query := fmt.Sprintf(`INSERT INTO post (%s) VALUES(LOWER($1), 1, $2, $3)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for insertPost: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(url_title, post.Title, posted)
	if err != nil {
		msg := "cannot execute query in insertPost: " + err.Error()
		return nil, errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in insertPost: " + err.Error()
		return nil, errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in insertPost but " + string(rows) + " rows were: " + err.Error()
		return nil, errors.New(msg)
	}
	post_id, err := res.LastInsertId()
	if err != nil {
		msg := "cannot get last insert id in insertPost: " + err.Error()
		return nil, errors.New(msg)
	}
	return &post_id, nil
}

func (database *Database) updatePost(tx *sqlx.Tx, db_post *Post, posted BOOL) error {
	query := `UPDATE post SET posted = $1, update_time = (CAST(strftime('%s', 'now') as integer)) WHERE id = $2`
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement for updatePost: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	res, err := stmt.Exec(posted, db_post.ID)
	if err != nil {
		msg := "cannot execute query in updatePost: " + err.Error()
		return errors.New(msg)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		msg := "cannot get affected rows in updatePost: " + err.Error()
		return errors.New(msg)
	}
	if rows != 1 {
		msg := "expected 1 row to be affected in updatePost but " + string(rows) + " rows were: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
