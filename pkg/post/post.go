package post

import (
	"errors"
	"regexp"

	"gitlab.com/joshraphael/diary/pkg/apierror"
	"gopkg.in/go-playground/validator.v9"
)

type Post struct {
	validator *validator.Validate
	method    string
	Title     string   `json:"title" validate:"required"`
	Tags      []string `json:"tags" validate:"required,min=1,max=10"`
	Body      string   `json:"body" validate:"required"`
}

func New(m string) Post {
	return Post{
		validator: validator.New(),
		method:    m,
	}
}

func (p Post) Method() string {
	return p.method
}

func (p Post) Validate() apierror.IApiError {
	if err := p.validator.Struct(p); err != nil {
		msg := "error validating post: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	urlSafe := regexp.MustCompile(`^[a-zA-Z0-9-_ ]{1,40}$`)
	if !urlSafe.MatchString(p.Title) {
		msg := "post title '" + p.Title + "' not URL safe"
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	for i := range p.Tags {
		if !urlSafe.MatchString(p.Tags[i]) {
			msg := "post tag '" + p.Tags[i] + "' not URL safe"
			apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
			return apiErr
		}
	}
	return nil
}
