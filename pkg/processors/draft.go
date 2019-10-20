package processors

import (
	"errors"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
)

func (prcr Processor) Draft(post_id int64, method string) (*database.CompletePost, apierror.IApiError) {
	db_post, err := prcr.db.GetPostById(post_id)
	if err != nil {
		msg := "error getting post in Draft: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", apierror.MethodHTTP)
		return nil, apiErr
	}
	if db_post == nil {
		msg := "No post exists for Draft"
		apiErr := apierror.New(errors.New(msg), "NOT_FOUND", apierror.MethodHTTP)
		return nil, apiErr
	}
	if db_post.Posted == database.DB_TRUE().Value() {
		msg := "Post has already been posted for Draft"
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", apierror.MethodHTTP)
		return nil, apiErr
	}
	post, err := prcr.db.GetCompletePost(db_post)
	if err != nil {
		msg := "cannot get complete posts: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", method)
		return nil, apiErr
	}
	if post == nil {
		msg := "No complete post exists for Draft"
		apiErr := apierror.New(errors.New(msg), "NOT_FOUND", apierror.MethodHTTP)
		return nil, apiErr
	}
	return post, nil
}
