package linter

import (
	"regexp"
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
func NewEndOfFileRule(config Config) Linter {
	return EndOfFileRule{
		Rule: Rule{
			Name:        "EOF Rule",
			Description: "Checks if file ends in a newline character.",
			AutoFix:     !config.Rules.EndOfFile.DisableAutofix && config.AutoFix,
			Ignore:      MustCompileIgnoreLines(append(config.Ignore, config.Rules.EndOfFile.Ignore...)...),
		},
		SingleNewLine: config.Rules.EndOfFile.SingleNewLine,
	}
}

// Lint implements the Lint interface
func (rule EndOfFileRule) Lint(b []byte) (valid bool, fix []byte) {

	// for empty files
	if len(b) < 1 {
		return true, nil
	}

	if rule.SingleNewLine && len(b) > 1 {
		valid = regexp.MustCompile(`[^\n]\n\z`).Match(b)
	} else {
		valid = regexp.MustCompile(`\n\z`).Match(b)
	}

	if valid || !rule.AutoFix {
		return valid, nil
	}

	// add one new line to the end of file
	fix = regexp.MustCompile(`(.)\z`).ReplaceAll(b, []byte("$1\n"))

	if rule.SingleNewLine {
		// rm extra new lines, if any
		fix = regexp.MustCompile(`\n+\z`).ReplaceAll(fix, []byte{'\n'})
	}

	return valid, fix
}
