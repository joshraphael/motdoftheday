package processors

import (
	"errors"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
)

func (prcr Processor) Edit(post_history_id int64, method string) (*database.CompletePostHistory, apierror.IApiError) {
	post_history, err := prcr.db.GetPostHistoryById(post_history_id)
	if err != nil {
		msg := "error getting post history in Edit: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", method)
		return nil, apiErr
	}
	post, err := prcr.db.GetPostById(post_history.PostID)
	if err != nil {
		msg := "error getting post in Edit: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", method)
		return nil, apiErr
	}
	if post.Posted == database.DB_TRUE().Value() {
		msg := "Post has already been posted for Edit"
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", apierror.MethodHTTP)
		return nil, apiErr
	}
	categories, err := prcr.db.GetPostHistoryCategories(post_history)
	if err != nil {
		msg := "error getting categories in Edit: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", method)
		return nil, apiErr
	}
	tags, err := prcr.db.GetPostHistoryTags(post_history)
	if err != nil {
		msg := "error getting tags in Edit: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", method)
		return nil, apiErr
	}
	return &database.CompletePostHistory{
		Post:       post,
		History:    post_history,
		Categories: categories,
		Tags:       tags,
	}, nil
}
