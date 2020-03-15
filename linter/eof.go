package linter

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"unicode/utf8"
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
			Name:        "New End of File",
			Description: "New End of File",
			Fix:         true,
			ignore:      setDefaultIgnore(),
		},
		SingleNewLine: true,
	}
}

// Lint lints
func (rule EndOfFileRule) Lint(r io.Reader) (valid bool, fix []byte, err error) {
	b, err := ioutil.ReadAll(r)

	if !isText(b) {
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

// isText reports whether a significant prefix of s looks like correct UTF-8;
// that is, if it is likely that s is human-readable text.
//
// see godoc:
// https://github.com/golang/tools/blob/gopls/v0.3.3/godoc/util/util.go#L38-L56
func isText(s []byte) bool {
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
