package main

import (
	"bufio"
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
	lint(r io.Reader) (valid bool, fix []byte, err error)
}

// Rule
type Rule struct {
	Name        string
	Description string
	Fix         bool
}

// NewLineRule
type NewLineRule struct {
	Rule
}

// SingleNewLineRule
type SingleNewLineRule struct {
	Rule
}

var newLineRule = NewLineRule{
	Rule{
		Name:        "New Line Rule",
		Description: "New Line Rule",
		Fix:         true,
	},
}

var singleNewLineRule = SingleNewLineRule{
	Rule{
		Name:        "Single New Line Rule",
		Description: "Single New Line Rule",
		Fix:         true,
	},
}

func (rule NewLineRule) lint(r io.Reader) (valid bool, fix []byte, err error) {
	b, err := ioutil.ReadAll(r)

	if !util.IsText(b) {
		return false, nil, fmt.Errorf("not text file")
	}

	if err != nil {
		return false, nil, nil
	}

	match := regexp.MustCompile(`[^\n]\n\z`).Match(b)

	if match || !rule.Fix {
		return match, nil, nil
	}

	// add one new line to the end of file
	fix = regexp.MustCompile(`(.)\z`).ReplaceAll(b, []byte("$1\n"))

	return match, fix, nil
}

func (rule SingleNewLineRule) lint(r io.Reader) (valid bool, fix []byte, err error) {
	b, err := ioutil.ReadAll(r)

	if !util.IsText(b) {
		return false, nil, fmt.Errorf("not text file")
	}

	if err != nil {
		return false, nil, nil
	}

	match := regexp.MustCompile(`[^\n]\n\z`).Match(b)

	if match || !rule.Fix {
		return match, nil, nil
	}

	// add one new line to the end of file
	fix = regexp.MustCompile(`(.)\z`).ReplaceAll(b, []byte("$1\n"))

	// rm extra new lines, if any
	fix = regexp.MustCompile(`\n+\z`).ReplaceAll(fix, []byte{'\n'})

	return match, fix, nil
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
		fr, err := os.Open(path)

		if err != nil {
			fmt.Printf("Error opening file %q: %v\n", path, err)
			os.Exit(1)
		}

		defer fr.Close()

		for _, rule := range rules {

			valid, fix, err := rule.lint(fr)

			if err != nil {
				fmt.Printf("Skipping file %q: %v\n", path, err)
			}

			if !valid {
				fmt.Printf("File %q has linting errors\n", path)
				errors++
			}

			// ignore errors
			fr.Close()

			if fix != nil {

				// will erase the file
				fw, err := os.Create(path)

				if err != nil {
					fmt.Printf("Failed to fix file %q: %v\n", path, err)
					break
				}

				defer fw.Close()

				w := bufio.NewWriter(fw)
				defer w.Flush()

				_, err = w.Write(fix)

				if err != nil {
					fmt.Printf("Failed to fix file %q: %v\n", path, err)
					break
				}

				err = w.Flush()

				if err != nil {
					fmt.Printf("Failed to flush file %q: %v\n", path, err)
					break
				}

				fmt.Printf("File %q linting errors fixed\n", path)
				errors--

				// ignore errors
				fw.Close()
			}
		}
	}

	if errors != 0 {
		fmt.Printf("\nTotal of %d linting errors!\n\n", errors)
		os.Exit(1)
	}
}
