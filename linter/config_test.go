package linter

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var yamlTestConfig = `
autofix: true

ignore:
  - .git/

rules:
  end-of-file:
    enable: true
    disable-autofix: false
    single-new-line: true
`

var defaultTestConf = newDefaultConfig()

func TestConfig(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal([]byte(yamlTestConfig), &c)
	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	if !reflect.DeepEqual(c, defaultTestConf) {
		t.Errorf("yaml.Unmarshal(Config):\n\tExpected %+v, got %+v", defaultTestConf, c)
	}
}
