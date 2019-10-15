package post

import (
	"errors"
	"regexp"
	"strings"

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

func (p Post) UrlTitle() string {
	return urlSafe(p.Title)
}

func (p Post) UrlTags() []string {
	tags := []string{}
	for i := range p.Tags {
		tags = append(tags, urlSafe(p.Tags[i]))
	}
	return tags
}

func (p Post) Validate() error {
	if err := p.validator.Struct(p); err != nil {
		msg := "error validating post: " + err.Error()
		return errors.New(msg)
	}
	urlSafe := regexp.MustCompile(`^[a-zA-Z0-9-_ ]{1,40}$`)
	if !urlSafe.MatchString(p.Title) {
		msg := "post title '" + p.Title + "' not URL safe"
		return errors.New(msg)
	}
	for i := range p.Tags {
		if !urlSafe.MatchString(p.Tags[i]) {
			msg := "post tag '" + p.Tags[i] + "' not URL safe"
			return errors.New(msg)
		}
	}
	return nil
}

func urlSafe(s string) string {
	return strings.Join(strings.Split(strings.TrimSpace(s), " "), "-")
}
