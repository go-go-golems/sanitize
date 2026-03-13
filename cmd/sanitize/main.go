package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("sanitize", flag.ContinueOnError)
	flags.SetOutput(stderr)

	jsonOut := flags.Bool("json", false, "output full result as JSON")
	lintOnly := flags.Bool("lint", false, "lint only (no fixing)")
	tabWidth := flags.Int("tab-width", 2, "spaces per tab for tab_indent fixer")
	maxIter := flags.Int("max-iterations", 10, "maximum fix iterations")

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	input, err := readInput(flags.Args(), stdin)
	if err != nil {
		fmt.Fprintf(stderr, "error reading input: %v\n", err)
		return 1
	}

	src := string(input)

	if *lintOnly {
		issues := yamlsanitize.Lint(src)
		if *jsonOut {
			if err := json.NewEncoder(stdout).Encode(issues); err != nil {
				fmt.Fprintf(stderr, "error encoding lint result: %v\n", err)
				return 1
			}
		} else {
			for _, li := range issues {
				fmt.Fprintln(stdout, li.Description)
			}
		}
		if len(issues) > 0 {
			return 1
		}
		return 0
	}

	result := yamlsanitize.Sanitize(src,
		yamlsanitize.WithTabWidth(*tabWidth),
		yamlsanitize.WithMaxIterations(*maxIter),
	)

	if *jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			fmt.Fprintf(stderr, "error encoding sanitize result: %v\n", err)
			return 1
		}
	} else {
		if _, err := io.WriteString(stdout, result.Sanitized); err != nil {
			fmt.Fprintf(stderr, "error writing sanitized output: %v\n", err)
			return 1
		}
		if len(result.Fixes) > 0 {
			fmt.Fprintf(stderr, "%d fix(es) applied\n", len(result.Fixes))
			for _, f := range result.Fixes {
				fmt.Fprintf(stderr, "  %s: %s\n", f.Rule, f.Description)
			}
		}
	}

	if !result.ParseClean || !result.LintClean {
		return 1
	}
	return 0
}

func readInput(args []string, stdin io.Reader) ([]byte, error) {
	if len(args) > 0 {
		return os.ReadFile(args[0])
	}
	return io.ReadAll(stdin)
}
