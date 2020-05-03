package processors

import "gitlab.com/joshraphael/motdoftheday/pkg/database"

type Processor struct {
	db  *database.Database
	cfg Config
}

func New(cfg Config, database *database.Database) Processor {
	return Processor{
		db:  database,
		cfg: cfg,
	}
}
