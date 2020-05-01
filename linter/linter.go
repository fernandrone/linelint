package linter

import (
	"unicode/utf8"

	gitignore "github.com/sabhiram/go-gitignore"
)

// MustCompileIgnoreLines compiles the ignore lines and throws a panic if it fails
func MustCompileIgnoreLines(lines ...string) *gitignore.GitIgnore {
	g, err := gitignore.CompileIgnoreLines(lines...)

	if err != nil {
		panic(err)
	}

	return g
}

// Linter exposes the lint method
type Linter interface {
	GetName() string

	// Lint performs the lint operation.
	//
	//  valid: 	if true, the file is valid, the linting check has passed
	//  fix:	if true, the linter will return a "fixed" copy of the file
	Lint([]byte) (valid bool, fix []byte)
	ShouldIgnore(path string) bool
}

// Rule represents a linting rule
type Rule struct {
	Config

	Name        string
	Description string

	// AutoFix sets if the linter should try to fix the error; if false, this Rule should
	// only return if the text is valid or not when the lint method is called; otherwise
	// the lint method should return a fixed copy of the data, to pass the linting rule
	AutoFix bool

	// ignore uses the gitignore syntax the select which files or folders to ignore
	Ignore *gitignore.GitIgnore
}

// GetName returns the rule name
func (r Rule) GetName() string {
	return r.Name
}

// ShouldIgnore decides if the file should be ignored based on the Ignore configuration
func (r Rule) ShouldIgnore(path string) bool {
	return r.Ignore.MatchesPath(path)
}

// IsText reports whether a significant prefix of s looks like correct UTF-8;
// that is, if it is likely that s is human-readable text.
//
// see godoc:
// https://github.com/golang/tools/blob/gopls/v0.3.3/godoc/util/util.go#L38-L56
func IsText(s []byte) bool {
	const max = 1024 // at least utf8.UTFMax
	if len(s) > max {
		s = s[0:max]
	}
	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			// last char may be incomplete - ignore
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			// decoding error or control character - not a text file
			return false
		}
	}
	return true
}
