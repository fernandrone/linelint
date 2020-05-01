package linter

import (
	"testing"

	"gopkg.in/yaml.v2"
)

var ignoreTests = []struct {
	file   string
	ignore bool
}{
	{"README", false},
	{".git/objects/04/9f2973ffc85f71da1fd5a", true},
}

func TestShouldIgnore_DefaultConf(t *testing.T) {
	for _, tt := range ignoreTests {
		t.Run(tt.file, func(t *testing.T) {
			got := NewEndOfFileRule(defaultTestConf).ShouldIgnore(tt.file)
			want := tt.ignore

			if got != want {
				t.Errorf("NewEndOfFileRule(defaultTestConf).ShouldIgnore(%q):\n\tExpected %v, got %v", tt.file, want, got)
			}
		})
	}
}

func TestShouldIgnore_YAMLParsedConf(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal([]byte(yamlTestConfig), &c)
	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	for _, tt := range ignoreTests {
		t.Run(tt.file, func(t *testing.T) {
			got := NewEndOfFileRule(c).ShouldIgnore(tt.file)
			want := tt.ignore

			if got != want {
				t.Errorf("NewEndOfFileRule(defaultTestConf).ShouldIgnore(%q):\n\tExpected %v, got %v", tt.file, want, got)
			}
		})
	}
}
