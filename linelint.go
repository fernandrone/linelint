package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/godoc/util"
)

// Linter exposes the lint method
type Linter interface {
	lint(r io.Reader) (valid bool, err error)
}

// Rule
type Rule struct {
	Name        string
	Description string
}

// NewLineRule
type NewLineRule struct {
	Rule
	NewLineRegex *regexp.Regexp
}

// SingleNewLineRule
type SingleNewLineRule struct {
	Rule
	SingleNewLineRegex *regexp.Regexp
}

var newLineRule = NewLineRule{
	Rule: Rule{
		Name:        "New Line Rule",
		Description: "New Line Rule",
	},
	NewLineRegex: regexp.MustCompile(`\n\z`),
}

var singleNewLineRule = SingleNewLineRule{
	Rule: Rule{
		Name:        "Single New Line Rule",
		Description: "Single New Line Rule",
	},
	SingleNewLineRegex: regexp.MustCompile(`[^\n]\n\z`),
}

func (rule NewLineRule) lint(r io.Reader) (valid bool, err error) {
	b, err := ioutil.ReadAll(r)

	if !util.IsText(b) {
		return false, fmt.Errorf("not text file")
	}

	if err != nil {
		return false, nil
	}

	return rule.NewLineRegex.Match(b), nil
}

func (rule SingleNewLineRule) lint(r io.Reader) (valid bool, err error) {
	b, err := ioutil.ReadAll(r)

	if !util.IsText(b) {
		return false, fmt.Errorf("not text file")
	}

	if err != nil {
		return false, nil
	}

	return rule.SingleNewLineRegex.Match(b), nil
}

func main() {
	var args, paths []string

	if len(os.Args[1:]) == 0 {
		args = []string{"."}
	} else {
		args = os.Args[1:]
	}

	for _, path := range args {
		f, err := os.Stat(path)

		if os.IsNotExist(err) {
			fmt.Printf("File %q does not exist", path)
			os.Exit(1)
		}

		// if dir, walk and append only files
		if f.IsDir() {
			err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", p, err)
					return err
				}

				// skip dirs
				if info.IsDir() {
					return nil
				}

				paths = append(paths, p)
				return nil
			})
			if err != nil {
				fmt.Printf("error walking the path %q: %v\n", path, err)
				return
			}
		} else {
			// if not dir, append
			paths = append(paths, path)
		}
	}

	var rules []Linter
	var errors int

	// for now we'll only use this rule, all the time
	rules = append(rules, singleNewLineRule)

	for _, path := range paths {
		file, err := os.Open(path)

		if err != nil {
			fmt.Printf("Error opening file %q: %e\n", path, err)
			os.Exit(1)
		}

		for _, rule := range rules {

			valid, err := rule.lint(file)

			if err != nil {
				fmt.Printf("Skipping file %q: %e\n", path, err)
			}

			if !valid {
				fmt.Printf("File %q has linting errors\n", path)
				errors++
			}
		}
	}

	if errors != 0 {
		fmt.Printf("\nTotal of %d linting errors!\n\n", errors)
		os.Exit(1)
	}
}
