package config

import (
	"gitlab.com/joshraphael/motdoftheday/internal/server/rest"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
	"gitlab.com/joshraphael/motdoftheday/pkg/processors"
)

type Config struct {
	Rest       rest.Config       `yaml:"rest" validate:"required"`
	Database   database.Config   `yaml:"db" validate:"required"`
	Processors processors.Config `yaml:"processors" validate:"required"`
}
