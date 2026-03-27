package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

func main() {
	root := flag.String("root", "examples/yaml", "Directory containing YAML example files")
	flag.Parse()

	paths, err := collectPaths(*root, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect paths: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("---")
	fmt.Println("Title: Parse vs Lint Matrix")
	fmt.Println("Ticket: SANITIZE-004")
	fmt.Println("Status: active")
	fmt.Println("Topics:")
	fmt.Println("    - yaml")
	fmt.Println("    - linting")
	fmt.Println("    - treesitter")
	fmt.Println("DocType: reference")
	fmt.Println("Intent: short-term")
	fmt.Println("Owners: []")
	fmt.Println("RelatedFiles: []")
	fmt.Println("ExternalSources: []")
	fmt.Println("Summary: \"Generated matrix comparing parse errors and lint issues across examples/yaml.\"")
	fmt.Printf("LastUpdated: %s\n", time.Now().Format(time.RFC3339Nano))
	fmt.Println("WhatFor: \"Generated evidence comparing parser and linter coverage across the YAML example corpus.\"")
	fmt.Println("WhenToUse: \"Use when validating tree-sitter-aware linting changes.\"")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# Parse vs Lint Matrix")
	fmt.Println()
	fmt.Printf("Generated from `%s` with `%d` file(s).\n", *root, len(paths))
	fmt.Println()
	fmt.Println("| File | Parse errors | Lint issues | Parse rows | Lint rows | Rules |")
	fmt.Println("|------|--------------|-------------|------------|-----------|-------|")

	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			os.Exit(1)
		}

		_, parseErrors, err := yamlsanitize.ParseTree(string(src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, err)
			os.Exit(1)
		}
		lintIssues := yamlsanitize.Lint(string(src))

		fmt.Printf(
			"| %s | %d | %d | %s | %s | %s |\n",
			path,
			len(parseErrors),
			len(lintIssues),
			formatParseRows(parseErrors),
			formatLintRows(lintIssues),
			formatRules(lintIssues),
		)
	}

	fmt.Println()
	fmt.Println("## Details")
	fmt.Println()

	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			os.Exit(1)
		}

		_, parseErrors, err := yamlsanitize.ParseTree(string(src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, err)
			os.Exit(1)
		}
		lintIssues := yamlsanitize.Lint(string(src))

		fmt.Printf("### %s\n\n", path)
		if len(parseErrors) == 0 {
			fmt.Println("- Parse errors: none")
		} else {
			fmt.Println("- Parse errors:")
			for _, parseErr := range parseErrors {
				fmt.Printf(
					"  - `%s` rows %d:%d to %d:%d text=%q\n",
					parseErr.Type,
					parseErr.StartRow+1,
					parseErr.StartCol+1,
					parseErr.EndRow+1,
					parseErr.EndCol+1,
					parseErr.Text,
				)
			}
		}

		if len(lintIssues) == 0 {
			fmt.Println("- Lint issues: none")
		} else {
			fmt.Println("- Lint issues:")
			for _, issue := range lintIssues {
				fmt.Printf("  - `%s` row %d: %s\n", issue.Rule, issue.Row+1, issue.Description)
			}
		}
		fmt.Println()
	}
}

func collectPaths(root string, args []string) ([]string, error) {
	if len(args) > 0 {
		paths := append([]string(nil), args...)
		sort.Strings(paths)
		return paths, nil
	}

	entries, err := filepath.Glob(filepath.Join(root, "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(entries)
	return entries, nil
}

func formatParseRows(errors []yamlsanitize.ErrorNode) string {
	if len(errors) == 0 {
		return "—"
	}
	rows := map[uint]struct{}{}
	for _, parseErr := range errors {
		rows[parseErr.StartRow+1] = struct{}{}
	}
	return joinUintRows(rows)
}

func formatLintRows(issues []yamlsanitize.LintIssue) string {
	if len(issues) == 0 {
		return "—"
	}
	rows := map[int]struct{}{}
	for _, issue := range issues {
		rows[issue.Row+1] = struct{}{}
	}
	return joinIntRows(rows)
}

func formatRules(issues []yamlsanitize.LintIssue) string {
	if len(issues) == 0 {
		return "—"
	}
	seen := map[string]struct{}{}
	rules := make([]string, 0, len(issues))
	for _, issue := range issues {
		if _, ok := seen[issue.Rule]; ok {
			continue
		}
		seen[issue.Rule] = struct{}{}
		rules = append(rules, issue.Rule)
	}
	sort.Strings(rules)
	return strings.Join(rules, ", ")
}

func joinUintRows(rows map[uint]struct{}) string {
	keys := make([]int, 0, len(rows))
	for row := range rows {
		keys = append(keys, int(row))
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, row := range keys {
		parts = append(parts, fmt.Sprintf("%d", row))
	}
	return strings.Join(parts, ", ")
}

func joinIntRows(rows map[int]struct{}) string {
	keys := make([]int, 0, len(rows))
	for row := range rows {
		keys = append(keys, row)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, row := range keys {
		parts = append(parts, fmt.Sprintf("%d", row))
	}
	return strings.Join(parts, ", ")
}
