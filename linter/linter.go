package linter

import (
	"io"

	gitignore "github.com/sabhiram/go-gitignore"
)

// setDefaultIgnore sets the default Ignore configuration using the gitignore syntax
func setDefaultIgnore() *gitignore.GitIgnore {
	return mustCompileIgnoreLines(".git")
}

// mustCompileIgnoreLines compiles the ignore lines and throws a panic if it fails
func mustCompileIgnoreLines(lines ...string) *gitignore.GitIgnore {
	g, err := gitignore.CompileIgnoreLines(lines...)

	if err != nil {
		panic(err)
	}

	return g
}

// Linter exposes the lint method
type Linter interface {
	Lint(r io.Reader) (valid bool, fix []byte, err error)
	ShouldIgnore(path string) bool
}

// Rule represents a linting rule
type Rule struct {
	Name        string
	Description string

	// Fix sets if the linter should try to fix the error; if false, this Rule should
	// only return if the text is valid or not when the lint method is called; otherwise
	// the lint method should return a fixed copy of the data, to pass the linting rule
	Fix bool

	// ignore uses the gitignore syntax the select which files or folders to ignore
	ignore *gitignore.GitIgnore
}

// ShouldIgnore uses decides if the file should be ignored based on the Ignore
// configuration
func (r Rule) ShouldIgnore(path string) bool {
	return r.ignore.MatchesPath(path)
}
