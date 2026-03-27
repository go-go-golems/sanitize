package main

import (
	"fmt"
	"sort"

	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	"sanitize007scripts/internal/corpus"
)

func main() {
	examples := corpus.LoadJSONExamples()
	buckets := map[string][]string{
		"clean":            {},
		"parser-driven":    {},
		"heuristic-driven": {},
		"hybrid":           {},
	}

	for _, ex := range examples {
		result := jsonsanitize.Sanitize(ex.Input)
		parserSignals := len(result.OriginalErrors)
		heuristicSignals := 0
		for _, issue := range result.OriginalLintIssues {
			if issue.Source == "heuristic" || issue.Source == "strict-parser" {
				heuristicSignals++
			}
		}

		bucket := "clean"
		switch {
		case parserSignals > 0 && heuristicSignals > 0:
			bucket = "hybrid"
		case parserSignals > 0:
			bucket = "parser-driven"
		case heuristicSignals > 0:
			bucket = "heuristic-driven"
		}

		label := ex.Name
		if ex.Filename != "" {
			label = ex.Filename + " — " + ex.Name
		}
		buckets[bucket] = append(buckets[bucket], label)
	}

	fmt.Println("---")
	fmt.Println("Title: JSON detection buckets")
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
	fmt.Println("Summary: Bucket classification showing which malformed JSON cases are parser-driven, heuristic-driven, or hybrid.")
	fmt.Println("WhatFor: Use when deciding where tree-sitter alone is enough and where heuristic lint/fix rules must remain first-class.")
	fmt.Println("WhenToUse: Read when prioritizing JSON rule work or evaluating recovery coverage.")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# JSON Detection Buckets")
	fmt.Println()

	order := []string{"clean", "parser-driven", "heuristic-driven", "hybrid"}
	for _, bucket := range order {
		sort.Strings(buckets[bucket])
		fmt.Printf("## %s\n\n", bucket)
		if len(buckets[bucket]) == 0 {
			fmt.Println("- none")
			fmt.Println()
			continue
		}
		for _, item := range buckets[bucket] {
			fmt.Printf("- %s\n", item)
		}
		fmt.Println()
	}
}
