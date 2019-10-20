package processors

import (
	"errors"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
)

func (prcr Processor) Drafts(method string) ([]database.Post, apierror.IApiError) {
	posts, err := prcr.db.GetDraftPosts()
	if err != nil {
		msg := "cannot get draft posts: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", method)
		return nil, apiErr
	}
	return posts, nil
}
