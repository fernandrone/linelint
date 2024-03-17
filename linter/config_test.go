package linter

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var yamlDefaultTestConfig = `
autofix: false

ignore:
  - .git/

rules:
  end-of-file:
    enable: true
    disable-autofix: false
    single-new-line: true
`

var yamlTestConfigWithIgnoreFile = `
autofix: false

ignore:
  - .git/
ignore-file: .gitignore

rules:
  end-of-file:
    enable: true
    disable-autofix: false
    single-new-line: true
`

var autofixTestConf = Config{
	AutoFix: true,
	Ignore:  []string{".git/"},
	Rules: RulesConfig{
		EndOfFile: EndOfFileConfig{
			Enable:        true,
			SingleNewLine: true,
		},
	},
}

var TestConfWithIgnoreFile = Config{
	AutoFix:    false,
	Ignore:     []string{".git/"},
	IgnoreFile: ".gitignore",
	Rules: RulesConfig{
		EndOfFile: EndOfFileConfig{
			Enable:        true,
			SingleNewLine: true,
		},
	},
}

func TestDefaultConfig(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal([]byte(yamlDefaultTestConfig), &c)
	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	if !reflect.DeepEqual(c, NewDefaultConfig()) {
		t.Errorf("yaml.Unmarshal(Config):\n\tExpected %+v, got %+v", NewDefaultConfig(), c)
	}
}

func TestConfigWithIgnoreFile(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal([]byte(yamlTestConfigWithIgnoreFile), &c)
	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	if !reflect.DeepEqual(c, TestConfWithIgnoreFile) {
		t.Errorf("yaml.Unmarshal(Config):\n\tExpected %+v, got %+v", TestConfWithIgnoreFile, c)
	}
}
