package processors

import "gitlab.com/joshraphael/motdoftheday/pkg/database"

type Processor struct {
	db *database.Database
}

func New(database *database.Database) Processor {
	return Processor{
		db: database,
	}
}
