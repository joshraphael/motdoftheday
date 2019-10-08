package rest

import "gopkg.in/go-playground/validator.v9"

type Rest struct {
	validator *validator.Validate
}

func New(v *validator.Validate) Rest {
	return Rest{
		validator: v,
	}
}
