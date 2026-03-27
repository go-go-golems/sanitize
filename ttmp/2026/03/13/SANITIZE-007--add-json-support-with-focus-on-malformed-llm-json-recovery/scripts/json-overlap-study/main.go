package main

import (
	"fmt"
	"sort"

	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	"sanitize007scripts/internal/corpus"
)

func main() {
	type summary struct {
		example    string
		total      int
		overlap    int
		rowOverlap int
	}

	summaries := make([]summary, 0)
	totalIssues := 0
	totalOverlap := 0
	totalRowOverlap := 0

	for _, ex := range corpus.LoadJSONExamples() {
		result := jsonsanitize.Sanitize(ex.Input)
		cur := summary{example: label(ex)}
		for _, issue := range result.OriginalLintIssues {
			if issue.Source != "heuristic" {
				continue
			}
			cur.total++
			totalIssues++
			if overlapsParse(issue, result.OriginalErrors) {
				cur.overlap++
				totalOverlap++
			}
			if overlapsParseRow(issue, result.OriginalErrors) {
				cur.rowOverlap++
				totalRowOverlap++
			}
		}
		summaries = append(summaries, cur)
	}

	sort.Slice(summaries, func(i, j int) bool { return summaries[i].example < summaries[j].example })

	fmt.Println("---")
	fmt.Println("Title: JSON heuristic overlap study")
	fmt.Println("Ticket: SANITIZE-007")
	fmt.Println("Status: active")
	fmt.Println("Topics:")
	fmt.Println("    - json")
	fmt.Println("    - linting")
	fmt.Println("DocType: reference")
	fmt.Println("Intent: long-term")
	fmt.Println("Owners: []")
	fmt.Println("RelatedFiles: []")
	fmt.Println("ExternalSources: []")
	fmt.Println("Summary: Overlap study comparing heuristic issue spans with tree-sitter parse-error spans across the JSON corpus.")
	fmt.Println("WhatFor: Show where parser spans reinforce heuristic confidence and where heuristics must stand on their own.")
	fmt.Println("WhenToUse: Use when deciding which heuristic fixes should remain automatic versus suggestion-only.")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# JSON Heuristic Overlap Study")
	fmt.Println()
	fmt.Printf("- Total heuristic issues: `%d`\n", totalIssues)
	fmt.Printf("- Byte-span overlap with parse errors: `%d`\n", totalOverlap)
	fmt.Printf("- Same-row overlap with parse errors: `%d`\n", totalRowOverlap)
	fmt.Println()
	fmt.Println("| Example | Heuristic issues | Byte overlap | Row overlap |")
	fmt.Println("| --- | --- | --- | --- |")
	for _, cur := range summaries {
		fmt.Printf("| %s | %d | %d | %d |\n", cur.example, cur.total, cur.overlap, cur.rowOverlap)
	}
	fmt.Println()
}

func label(ex corpus.Example) string {
	if ex.Filename != "" {
		return ex.Filename
	}
	return ex.Name
}

func overlapsParse(issue jsonsanitize.LintIssue, errors []jsonsanitize.ErrorNode) bool {
	for _, parseErr := range errors {
		if issue.StartByte < parseErr.EndByte && parseErr.StartByte < issue.EndByte {
			return true
		}
	}
	return false
}

func overlapsParseRow(issue jsonsanitize.LintIssue, errors []jsonsanitize.ErrorNode) bool {
	for _, parseErr := range errors {
		if issue.StartRow >= parseErr.StartRow && issue.StartRow <= parseErr.EndRow {
			return true
		}
	}
	return false
}
