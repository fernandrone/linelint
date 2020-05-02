package linter

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents a Rule configuration
type Config struct {
	// AutoFix sets if the linter should try to fix the error
	AutoFix bool `yaml:"autofix"`

	// Ignore uses the gitignore syntax the select which files or folders to ignore
	Ignore []string `yaml:"ignore"`

	// Rules contains configuration specific for each rule
	Rules RulesConfig `yaml:"rules"`
}

// RulesConfig is the map of configs which includes all rules
type RulesConfig struct {
	EndOfFile EndOfFileConfig `yaml:"end-of-file"`
}

// EndOfFileConfig config for the End of File rule
type EndOfFileConfig struct {
	Enable bool `yaml:"enable"`

	DisableAutofix bool `yaml:"disable-autofix"`

	// Ignore uses the gitignore syntax the select which files or folders to ignore
	Ignore []string `yaml:"ignore"`

	SingleNewLine bool `yaml:"single-new-line"`
}

// NewConfig returns a new Config
func NewConfig() Config {
	path := ".linelint.yml"

	var data []byte

	// check if config file exists
	if _, err := os.Stat(path); err != nil {
		return newDefaultConfig()
	}

	// if config file does exist, read it
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading YAML file %s: %s (will use default configuration)\n", path, err)
		return newDefaultConfig()
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing YAML file: %s (will use default configuration)\n", err)
		return newDefaultConfig()
	}

	return config
}

func newDefaultConfig() Config {
	return Config{
		AutoFix: true,
		Ignore:  []string{".git/"},
		Rules: RulesConfig{
			EndOfFile: EndOfFileConfig{
				Enable:        true,
				SingleNewLine: true,
			},
		},
	}
}
