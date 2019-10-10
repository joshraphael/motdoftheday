package processors

import (
	"errors"

	"gitlab.com/joshraphael/diary/pkg/apierror"
	"gitlab.com/joshraphael/diary/pkg/post"
)

func (prcr Processor) SaveForm(p post.Post) apierror.IApiError {
	err := p.Validate()
	if err != nil {
		msg := "invalid save post: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	return nil
}
