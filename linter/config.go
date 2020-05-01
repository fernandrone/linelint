package linter

// Config represents a Rule configuration
type Config struct {
	// Files represents a list of files that the rule will parse
	Files []string `yaml:"files"`

	// Ignores uses the gitignore syntax the select which files or folders to ignore
	Ignore string `yaml:"ignore"`
}

// NewConfig returns a new Config
// For now we only support a default configuration that includes all files but ignores
// the .git folder
func NewConfig() Config {
	return Config{
		Files:  []string{"*"},
		Ignore: ".git/",
	}
}
