package linter

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestShouldIgnore_DefaultConf(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal([]byte(`ignore: [ ".git/" ]`), &c)

	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	ignoreTests := []struct {
		file   string
		ignore bool
	}{
		{".git/objects/04/9f2973ffc85f71da1fd5a", true},
		{"README", false},
		{"clusters/mycluster/applications/app.yml", false},
		{"java/bin/myclass.class", false},
	}

	for _, tt := range ignoreTests {
		t.Run(tt.file, func(t *testing.T) {
			got := NewEndOfFileRule(c).ShouldIgnore(tt.file)
			want := tt.ignore

			if got != want {
				t.Errorf(
					"NewEndOfFileRule(c).ShouldIgnore(%q):\n\tExpected %v, got %v",
					tt.file, want, got,
				)
			}
		})
	}
}

func TestShouldIgnore_MoreComplexConf(t *testing.T) {
	c := Config{}

	err := yaml.Unmarshal(
		[]byte(`ignore: [ ".git/", "**/bin/", "applications", "/projects/" ]`), &c,
	)

	if err != nil {
		t.Fatalf("yaml.Unmarshal(Config): %v", err)
	}

	ignoreTests := []struct {
		file   string
		ignore bool
	}{
		{".git/objects/04/9f2973ffc85f71da1fd5a", true},
		{"README", false},
		{"clusters/mycluster/applications/app.yml", true},
		{"home/projects/data.md", false},
		{"projects/data.md", true},
		{"java/bin/myclass.class", true},
	}

	for _, tt := range ignoreTests {
		t.Run(tt.file, func(t *testing.T) {
			got := NewEndOfFileRule(c).ShouldIgnore(tt.file)
			want := tt.ignore

			if got != want {
				t.Errorf(
					"NewEndOfFileRule(c).ShouldIgnore(%q):\n\tExpected %v, got %v",
					tt.file, want, got,
				)
			}
		})
	}
}
