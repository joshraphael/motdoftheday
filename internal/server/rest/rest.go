package rest

import (
	"gitlab.com/joshraphael/motdoftheday/pkg/processors"
	"gopkg.in/go-playground/validator.v9"
)

type Rest struct {
	validator *validator.Validate
	processor processors.Processor
}

func New(v *validator.Validate, p processors.Processor) Rest {
	return Rest{
		validator: v,
		processor: p,
	}
}
