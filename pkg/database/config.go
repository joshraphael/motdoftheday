package database

type Config struct {
	File string `yaml:"file" validate:"required,min=1"`
}
