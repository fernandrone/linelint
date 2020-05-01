package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fernandrone/linelint/linter"
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

	var fileErrors, lintErrors int
	var linters []linter.Linter

	// for now we'll only use this rule, all the time
	linters = append(linters, linter.NewEndOfFileRule())

	for _, path := range paths {

		fr, err := os.Open(path)

		if err != nil {
			fmt.Printf("Error opening file %q: %v\n", path, err)
			fileErrors++
			continue
		}

		defer fr.Close()

		if err != nil {
			fmt.Printf("Skipping file %q: %v\n", path, err)
			continue
		}

		b, err := ioutil.ReadAll(fr)

		if err != nil {
			fmt.Printf("Error reading file %q: %v\n", path, err)
			fileErrors++
			continue
		}

		if !linter.IsText(b) {
			fmt.Printf("Ignoring file %q: not text file\n", path)
			continue
		}

		for _, rule := range linters {

			if rule.ShouldIgnore(path) {
				fmt.Printf("Ignoring file %q: in ignore path\n", path)
				continue
			}

			valid, fix := rule.Lint(b)

			if !valid {
				fmt.Printf("File %q has lint errors\n", path)
				lintErrors++
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

				fmt.Printf("File %q lint errors fixed\n", path)
				lintErrors--

				// ignore errors
				fw.Close()
			}
		}
	}

	if fileErrors != 0 {
		fmt.Printf("\nTotal of %d file errors!\n\n", fileErrors)
	}

	if lintErrors != 0 {
		fmt.Printf("\nTotal of %d lint errors!\n\n", lintErrors)
	}

	if fileErrors != 0 || lintErrors != 0 {
		os.Exit(1)
	}
}
