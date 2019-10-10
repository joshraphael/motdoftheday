package processors

import (
	"errors"

	"gitlab.com/joshraphael/diary/pkg/apierror"
	"gitlab.com/joshraphael/diary/pkg/post"
)

func (prcr Processor) SubmitForm(p post.Post) apierror.IApiError {
	err := p.Validate()
	if err != nil {
		msg := "invalid submit post: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	return nil
}
