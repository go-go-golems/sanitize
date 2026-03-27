package main

import (
	"fmt"
	"sort"
	"strings"

	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	"sanitize007scripts/internal/corpus"
)

func main() {
	fmt.Println("---")
	fmt.Println("Title: JSON repair matrix")
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
	fmt.Println("Summary: End-to-end repair matrix for JSON examples, including parser signals, lint rules, fixes, and final strict-parse status.")
	fmt.Println("WhatFor: Show which malformed JSON cases are recoverable today and which remain lint-only.")
	fmt.Println("WhenToUse: Use when reviewing JSON recovery coverage or selecting the next rule/fixer to implement.")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# JSON Repair Matrix")
	fmt.Println()
	fmt.Println("| Example | Pattern | Original strict | Parse nodes | Rules | Fixes | Final strict |")
	fmt.Println("| --- | --- | --- | --- | --- | --- | --- |")

	for _, ex := range corpus.LoadJSONExamples() {
		result := jsonsanitize.Sanitize(ex.Input)
		fmt.Printf(
			"| %s | `%s` | `%t` | `%d` | `%s` | `%s` | `%t` |\n",
			label(ex),
			ex.Pattern,
			result.OriginalStrictParseClean,
			len(result.OriginalErrors),
			strings.Join(uniqueRuleNames(result.OriginalLintIssues), ", "),
			strings.Join(uniqueFixNames(result.Fixes), ", "),
			result.StrictParseClean,
		)
	}

	fmt.Println()
	for _, ex := range corpus.LoadJSONExamples() {
		result := jsonsanitize.Sanitize(ex.Input)
		fmt.Printf("## %s\n\n", label(ex))
		fmt.Printf("- Pattern: `%s`\n", ex.Pattern)
		fmt.Printf("- Original strict parse clean: `%t`\n", result.OriginalStrictParseClean)
		fmt.Printf("- Original tree-sitter error count: `%d`\n", len(result.OriginalErrors))
		fmt.Printf("- Original rules: `%s`\n", strings.Join(uniqueRuleNames(result.OriginalLintIssues), ", "))
		fmt.Printf("- Applied fixes: `%s`\n", strings.Join(uniqueFixNames(result.Fixes), ", "))
		fmt.Printf("- Final strict parse clean: `%t`\n", result.StrictParseClean)
		fmt.Printf("- Final lint clean: `%t`\n", result.LintClean)
		fmt.Println()
		fmt.Println("```text")
		fmt.Println(strings.TrimSuffix(result.Sanitized, "\n"))
		fmt.Println("```")
		fmt.Println()
	}
}

func label(ex corpus.Example) string {
	if ex.Filename != "" {
		return ex.Filename
	}
	return ex.Name
}

func uniqueRuleNames(issues []jsonsanitize.LintIssue) []string {
	set := map[string]struct{}{}
	for _, issue := range issues {
		set[issue.Rule] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for rule := range set {
		out = append(out, rule)
	}
	sort.Strings(out)
	if len(out) == 0 {
		return []string{"none"}
	}
	return out
}

func uniqueFixNames(fixes []jsonsanitize.Fix) []string {
	set := map[string]struct{}{}
	for _, fix := range fixes {
		set[fix.Rule] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for rule := range set {
		out = append(out, rule)
	}
	sort.Strings(out)
	if len(out) == 0 {
		return []string{"none"}
	}
	return out
}
