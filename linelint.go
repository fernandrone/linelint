package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

// matches a single newline
var newline = regexp.MustCompile(`[^\n]\n\z`)

func main() {
	var paths []string

	for _, path := range os.Args[1:] {
		f, err := os.Stat(path)

		if os.IsNotExist(err) {
			fmt.Printf("File %q does not exist", path)
			os.Exit(1)
		}

		if f.IsDir() {
			fmt.Printf("Path %q is a directory, not a file", path)
		}

		paths = append(paths, path)
	}

	var errors int

	for _, path := range paths {
		file, err := os.Open(path)

		if err != nil {
			fmt.Printf("Error opening file %q: %e\n", path, err)
			os.Exit(1)
		}

		valid, err := lintFile(file)

		if err != nil {
			fmt.Printf("Error reading file %q: %e\n", path, err)
			os.Exit(1)
		}

		if !valid {
			fmt.Printf("File %q has linting errors\n", path)
			errors++
		}
	}

	if errors != 0 {
		fmt.Printf("\nTotal of %d linting errors!\n\n", errors)
		os.Exit(1)
	}
}

func lintFile(f io.Reader) (bool, error) {
	b, err := ioutil.ReadAll(f)

	if err != nil {
		return false, err
	}

	return lintBytes(b), err
}

func lintBytes(b []byte) bool {
	return newline.Match(b)
}
