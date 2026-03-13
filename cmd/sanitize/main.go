package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

func main() {
	jsonOut := flag.Bool("json", false, "output full result as JSON")
	lintOnly := flag.Bool("lint", false, "lint only (no fixing)")
	tabWidth := flag.Int("tab-width", 2, "spaces per tab for tab_indent fixer")
	maxIter := flag.Int("max-iterations", 10, "maximum fix iterations")
	flag.Parse()

	var input []byte
	var err error

	if flag.NArg() > 0 {
		input, err = os.ReadFile(flag.Arg(0))
	} else {
		input, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}

	src := string(input)

	if *lintOnly {
		issues := yamlsanitize.Lint(src)
		if *jsonOut {
			_ = json.NewEncoder(os.Stdout).Encode(issues)
		} else {
			for _, li := range issues {
				fmt.Println(li.Description)
			}
			if len(issues) > 0 {
				os.Exit(1)
			}
		}
		return
	}

	result := yamlsanitize.Sanitize(src,
		yamlsanitize.WithTabWidth(*tabWidth),
		yamlsanitize.WithMaxIterations(*maxIter),
	)

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
	} else {
		fmt.Print(result.Sanitized)
		if len(result.Fixes) > 0 {
			fmt.Fprintf(os.Stderr, "%d fix(es) applied\n", len(result.Fixes))
			for _, f := range result.Fixes {
				fmt.Fprintf(os.Stderr, "  %s: %s\n", f.Rule, f.Description)
			}
		}
		if !result.ParseClean || !result.LintClean {
			os.Exit(1)
		}
	}
}
