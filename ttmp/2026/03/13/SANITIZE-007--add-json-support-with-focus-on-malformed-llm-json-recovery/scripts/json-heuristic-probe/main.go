package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type testCase struct {
	Number       string
	Title        string
	Language     string
	Snippet      string
	ExpectedNote string
}

func main() {
	input := filepath.Join("..", "sources", "local", "03-json-parse-errors-import.md")
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	content, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input: %v\n", err)
		os.Exit(1)
	}

	cases := parseCases(string(content))
	fmt.Println("---")
	fmt.Println("Title: JSON heuristic probe")
	fmt.Println("Ticket: SANITIZE-007")
	fmt.Println("Status: active")
	fmt.Println("Topics:")
	fmt.Println("    - json")
	fmt.Println("    - linting")
	fmt.Println("    - api-design")
	fmt.Println("DocType: reference")
	fmt.Println("Intent: long-term")
	fmt.Println("Owners: []")
	fmt.Println("RelatedFiles: []")
	fmt.Println("ExternalSources: []")
	fmt.Println("Summary: Generated heuristic hit report showing which malformed JSON cases are detectable by simple string/regex heuristics before structural parsing.")
	fmt.Println("WhatFor: Show where simple heuristics are enough and where tree-sitter or stricter parsing is still required.")
	fmt.Println("WhenToUse: Use when deciding which malformed JSON cases should become heuristic lint rules.")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# JSON Heuristic Probe")
	fmt.Println()
	fmt.Printf("Source: `%s`\n\n", input)
	for _, tc := range cases {
		fmt.Printf("## Case %s: %s\n\n", tc.Number, tc.Title)
		if tc.Snippet == "" {
			fmt.Println("- Snippet: none in source note")
			fmt.Println()
			continue
		}
		for _, hit := range detectHeuristics(tc.Snippet) {
			fmt.Printf("- `%s`\n", hit)
		}
		fmt.Println()
	}
}

func parseCases(content string) []testCase {
	headerRe := regexp.MustCompile(`^(\d+)\.\s+\*\*(.+?)\*\*`)
	lines := strings.Split(content, "\n")
	var cases []testCase
	var current *testCase

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if m := headerRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				cases = append(cases, *current)
			}
			current = &testCase{Number: m[1], Title: m[2]}
			continue
		}
		if current == nil {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Error:") {
			current.ExpectedNote = strings.TrimPrefix(trimmed, "Error: ")
			continue
		}

		if fence, lang, ok := fenceParts(trimmed); ok {
			current.Language = lang
			var block []string
			for i = i + 1; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == fence {
					break
				}
				block = append(block, lines[i])
			}
			current.Snippet = strings.Join(block, "\n")
		}
	}
	if current != nil {
		cases = append(cases, *current)
	}
	return cases
}

func fenceParts(trimmed string) (string, string, bool) {
	if !strings.HasPrefix(trimmed, "```") {
		return "", "", false
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == '`' {
		n++
	}
	if n < 3 {
		return "", "", false
	}
	return strings.Repeat("`", n), strings.TrimSpace(trimmed[n:]), true
}

func detectHeuristics(snippet string) []string {
	hits := []string{}
	add := func(name string, ok bool) {
		if ok {
			hits = append(hits, name)
		}
	}

	add("markdown_fence", strings.Contains(snippet, "```"))
	add("comment", strings.Contains(snippet, "//") || strings.Contains(snippet, "/*"))
	add("single_quotes", strings.Contains(snippet, "'"))
	add("python_literals", strings.Contains(snippet, "True") || strings.Contains(snippet, "False") || strings.Contains(snippet, "None"))
	add("trailing_comma", regexp.MustCompile(`,\s*[}\]]`).MatchString(snippet))
	add("duplicate_comma", strings.Contains(snippet, ",,"))
	add("ellipsis_or_placeholder", strings.Contains(snippet, "...") || strings.Contains(snippet, "<") || strings.Contains(snippet, "[insert"))
	add("prose_wrapped", strings.Contains(snippet, "Here is your JSON:") || strings.Contains(snippet, "Thanks!"))
	add("multiple_top_level", strings.Contains(snippet, "} {"))
	add("unquoted_keys", regexp.MustCompile(`(?m)[{,]\s*[A-Za-z_][A-Za-z0-9_]*\s*:`).MatchString(snippet))

	if len(hits) == 0 {
		hits = append(hits, "none")
	}
	return hits
}
