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
	Verbose bool `yaml:"verbose"`

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

// NewConfigFromFile returns a new Config
func NewConfigFromFile(path string) Config {
	var data []byte

	// check if config file exists
	if _, err := os.Stat(path); err != nil {
		return NewDefaultConfig()
	}

	// if config file does exist, read it
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading YAML file %s: %s (will use default configuration)\n", path, err)
		return NewDefaultConfig()
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing YAML file: %s (will use default configuration)\n", err)
		return NewDefaultConfig()
	}

	return config
}

func NewDefaultConfig() Config {
	return Config{
		AutoFix: false,
		Ignore:  []string{".git/"},
		Rules: RulesConfig{
			EndOfFile: EndOfFileConfig{
				Enable:        true,
				SingleNewLine: true,
			},
		},
	}
}
