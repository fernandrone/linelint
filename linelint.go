package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fernandrone/linelint/rules"
)

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

	var r []rules.Linter
	var errors int

	// for now we'll only use this rule, all the time
	r = append(r, rules.NewEndOfFileRule())

	for _, path := range paths {
		fr, err := os.Open(path)

		if err != nil {
			fmt.Printf("Error opening file %q: %v\n", path, err)
			os.Exit(1)
		}

		defer fr.Close()

		for _, rule := range r {

			valid, fix, err := rule.Lint(fr)

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
