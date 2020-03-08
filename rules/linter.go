package rules

import "io"

// Linter exposes the lint method
type Linter interface {
	Lint(r io.Reader) (valid bool, fix []byte, err error)
}

// Rule represents a linting rule
type Rule struct {
	Name        string
	Description string

	// Fix sets if the linter should try to fix the error; if false, this Rule should
	// only return if the text is valid or not when the lint method is called; otherwise
	// the lint method should return a fixed copy of the data, to pass the linting rule
	Fix bool
}
