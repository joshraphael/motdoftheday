package processors

type Config struct {
	Directory    string `yaml:"dir" validate:"required"`
	TemplateFile string `yaml:"template" validate:"required"`
}
