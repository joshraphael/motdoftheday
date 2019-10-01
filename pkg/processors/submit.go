package processors

import (
	"errors"
	"time"
)

type Submit struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Body  string    `json:"body"`
}

func SubmitForm(submit Submit) error {
	if len(submit.Title) == 0 {
		msg := "Error missing title in json data"
		return errors.New(msg)
	}
	if submit.Date.IsZero() {
		msg := "Error missing date in json data"
		return errors.New(msg)
	}
	if len(submit.Body) == 0 {
		msg := "Error missing body in json data"
		return errors.New(msg)
	}
	return nil
}
