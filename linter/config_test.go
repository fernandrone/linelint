package linter

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var config = `
files:
  - "*"

ignore: |
  .git/`

func TestConfig(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal([]byte(config), &c)
	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	files := []string{"*"}

	if !reflect.DeepEqual(c.Files, files) {
		t.Errorf("yaml.Unmarshal(Config).Files:\n\tExpected %q, got %q", files, c.Files)
	}

	ignore := ".git/"

	if c.Ignore != ignore {
		t.Errorf("yaml.Unmarshal(Config).Ignore:\n\tExpected %q, got %q", ignore, c.Ignore)
	}
}
