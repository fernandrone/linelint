package rules

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"golang.org/x/tools/godoc/util"
)

// EndOfFileRule checks if the file ends in a newline character, or `\n`. It can be
// configured to check if it ends strictly in a single newline or in any number of
// newline characters.
type EndOfFileRule struct {
	Rule

	// If SingleNewLine is true, the EndOfFileRule expects that files end strictly in a
	// single newline character.
	SingleNewLine bool
}

// NewEndOfFileRule returns a new EndOfFileRule
func NewEndOfFileRule() Linter {
	return EndOfFileRule{
		Rule: Rule{
			Name:        "New Line Rule",
			Description: "New Line Rule",
			Fix:         true,
		},
		SingleNewLine: true,
	}
}

// Lint lints
func (rule EndOfFileRule) Lint(r io.Reader) (valid bool, fix []byte, err error) {
	b, err := ioutil.ReadAll(r)

	if !util.IsText(b) {
		return false, nil, fmt.Errorf("not text file")
	}

	if err != nil {
		return false, nil, nil
	}

	if rule.SingleNewLine {
		valid = regexp.MustCompile(`[^\n]\n\z`).Match(b)
	} else {
		valid = regexp.MustCompile(`\n\z`).Match(b)
	}

	if valid || !rule.Fix {
		return valid, nil, nil
	}

	// add one new line to the end of file
	fix = regexp.MustCompile(`(.)\z`).ReplaceAll(b, []byte("$1\n"))

	if rule.SingleNewLine {
		// rm extra new lines, if any
		fix = regexp.MustCompile(`\n+\z`).ReplaceAll(fix, []byte{'\n'})
	}

	return valid, fix, nil
}
