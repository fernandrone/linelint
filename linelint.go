package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fernandrone/linelint/linter"
)

const helpMsg = `usage of %s [-a] [FILE_OR_DIR [FILE_OR_DIR ...]]

Validates simple newline and whitespace rules in all sorts of files.

positional arguments:
  FILE_OR_DIR		files to format or '-' for stdin

optional arguments:
`

// Input is the main input structure to the program
type Input struct {
	Paths  []string
	Stdin  io.Reader
	Config linter.Config
}

func main() {
	var flagAutofix bool
	flag.BoolVar(&flagAutofix, "a", false, "(autofix) will automatically fix files with errors in place")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpMsg, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var paths []string

	if flag.NArg() == 0 {
		paths = []string{"."}
	} else {
		paths = flag.Args()
	}

	config := linter.NewConfig()

	if flagAutofix {
		config.AutoFix = true
	}

	input := Input{
		Paths:  paths,
		Stdin:  os.Stdin,
		Config: config,
	}

	if err := run(input); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(in Input) error {
	var linters []linter.Linter

	if in.Config.Rules.EndOfFile.Enable {
		linters = append(linters, linter.NewEndOfFileRule(in.Config))
	}

	if len(linters) == 0 {
		return errors.New("No valid rule enabled")
	}

	// read from stdin
	if len(in.Paths) == 1 && in.Paths[0] == "-" {
		return processSTDIN(in, linters)
	}

	return processDirectoryTree(in, linters)
}

func processSTDIN(in Input, linters []linter.Linter) error {
	var lintErrors int

	b, err := ioutil.ReadAll(in.Stdin)

	if err != nil {
		return fmt.Errorf("Error reading from Stdin: %v", err)
	}

	if !linter.IsText(b) {
		return errors.New("Stdin is not a valid UFT-8 input")
	}

	for _, rule := range linters {

		valid, fix := rule.Lint(b)

		if !valid {
			lintErrors++
		}

		if fix != nil {

			if err != nil {
				return fmt.Errorf("[%s] Failed to fix Stdin: %v\n", rule.GetName(), err)
			}

			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			_, err = w.Write(fix)

			if err != nil {
				return fmt.Errorf("[%s] Failed to print fixed input to Stdout: %v\n", rule.GetName(), err)
			}

			err = w.Flush()

			if err != nil {
				return fmt.Errorf("[%s] Failed to flush fixed input to Stdout: %v\n", rule.GetName(), err)
			}

			lintErrors--
		}
	}

	if lintErrors != 0 {
		// call exit directly to disable the error message
		os.Exit(1)
	}

	return nil
}

func processDirectoryTree(in Input, linters []linter.Linter) error {
	var files []string

	// get patterns to ignore
	ignore := linter.MustCompileIgnoreLines(in.Config.Ignore...)

	for _, path := range in.Paths {
		f, err := os.Stat(path)

		if os.IsNotExist(err) {
			return fmt.Errorf("File %q does not exist", path)
		}

		// if dir, walk and append only files
		if f.IsDir() {
			err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("Prevent panic by handling failure accessing a path %q: %v\n", p, err)
					return err
				}

				// skip dirs
				if info.IsDir() {
					return nil
				}

				if ignore.MatchesPath(p) {
					return nil
				}

				files = append(files, p)
				return nil
			})
			if err != nil {
				return fmt.Errorf("Error walking the path %q: %v\n", path, err)
			}
		} else {
			// if not dir, append
			files = append(files, path)
		}
	}

	var fileErrors, lintErrors int

	for _, f := range files {

		fr, err := os.Open(f)

		if err != nil {
			fmt.Printf("Error opening file %q: %v\n", f, err)
			fileErrors++
			continue
		}

		defer fr.Close()

		if err != nil {
			fmt.Printf("Skipping file %q: %v\n", f, err)
			continue
		}

		b, err := ioutil.ReadAll(fr)

		if err != nil {
			fmt.Printf("Error reading file %q: %v\n", f, err)
			fileErrors++
			continue
		}

		if !linter.IsText(b) {
			// TODO add log levels
			// fmt.Printf("Ignoring file %q: not text file\n", path)
			continue
		}

		for _, rule := range linters {

			if rule.ShouldIgnore(f) {
				fmt.Printf("[%s] Ignoring file %q: in rule ignore path\n", rule.GetName(), f)
				continue
			}

			valid, fix := rule.Lint(b)

			if !valid {
				fmt.Printf("[%s] File %q has lint errors\n", rule.GetName(), f)
				lintErrors++
			}

			// ignore errors
			fr.Close()

			if fix != nil {
				// will erase the file
				fw, err := os.Create(f)

				if err != nil {
					fmt.Printf("[%s] Failed to fix file %q: %v\n", rule.GetName(), f, err)
					break
				}

				defer fw.Close()

				w := bufio.NewWriter(fw)
				defer w.Flush()

				_, err = w.Write(fix)

				if err != nil {
					fmt.Printf("[%s] Failed to fix file %q: %v\n", rule.GetName(), f, err)
					break
				}

				err = w.Flush()

				if err != nil {
					fmt.Printf("[%s] Failed to flush file %q: %v\n", rule.GetName(), f, err)
					break
				}

				fmt.Printf("[%s] File %q lint errors fixed\n", rule.GetName(), f)
				lintErrors--

				// ignore errors
				fw.Close()
			}
		}
	}

	if fileErrors != 0 {
		fmt.Printf("\nTotal of %d file errors!\n", fileErrors)
	}

	if lintErrors != 0 {
		fmt.Printf("\nTotal of %d lint errors!\n", lintErrors)
	}

	if fileErrors != 0 || lintErrors != 0 {
		// call exit directly to disable the error message
		os.Exit(1)
	}

	return nil
}
