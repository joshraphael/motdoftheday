package processors

import (
	"errors"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

func (prcr Processor) SubmitForm(p post.Post) apierror.IApiError {
	err := p.Validate()
	if err != nil {
		msg := "invalid submit post: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	err = prcr.db.CreatePost(p, database.DB_TRUE())
	if err != nil {
		msg := "cannot submit post: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	return nil
}
