package main

import (
	stdjson "encoding/json"
	"fmt"
	"os"
	"sort"

	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	"sanitize007scripts/internal/corpus"
)

type row struct {
	Example             string   `json:"example"`
	Pattern             string   `json:"pattern"`
	PrimaryRule         string   `json:"primary_rule"`
	Severity            string   `json:"severity"`
	FixConfidence       string   `json:"fix_confidence"`
	AutoFixPolicy       string   `json:"auto_fix_policy"`
	ParserSignal        string   `json:"parser_signal"`
	ObservedRules       []string `json:"observed_rules"`
	ObservedFixes       []string `json:"observed_fixes"`
	OriginalStrictClean bool     `json:"original_strict_clean"`
	FinalStrictClean    bool     `json:"final_strict_clean"`
}

func main() {
	examples := corpus.LoadJSONExamples()
	rows := make([]row, 0, len(examples))
	for _, ex := range examples {
		result := jsonsanitize.Sanitize(ex.Input)
		meta := policyForPattern(ex.Pattern)
		rows = append(rows, row{
			Example:             label(ex),
			Pattern:             ex.Pattern,
			PrimaryRule:         meta.PrimaryRule,
			Severity:            meta.Severity,
			FixConfidence:       meta.FixConfidence,
			AutoFixPolicy:       meta.AutoFixPolicy,
			ParserSignal:        parserSignal(result),
			ObservedRules:       uniqueRuleNames(result.OriginalLintIssues),
			ObservedFixes:       uniqueFixNames(result.Fixes),
			OriginalStrictClean: result.OriginalStrictParseClean,
			FinalStrictClean:    result.StrictParseClean,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Example < rows[j].Example
	})

	enc := stdjson.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rows); err != nil {
		fmt.Fprintf(os.Stderr, "encode rule matrix: %v\n", err)
		os.Exit(1)
	}
}

type policy struct {
	PrimaryRule   string
	Severity      string
	FixConfidence string
	AutoFixPolicy string
}

func policyForPattern(pattern string) policy {
	switch pattern {
	case "markdown_fence_wrapper", "leading_or_trailing_prose", "python_literals", "duplicate_comma", "comment", "trailing_comma":
		return policy{PrimaryRule: pattern, Severity: "warn", FixConfidence: "safe-auto-fix", AutoFixPolicy: "auto-fix"}
	case "wrapper_python_trailing_comma_combo", "comment_python_duplicate_comma_combo":
		return policy{PrimaryRule: pattern, Severity: "error", FixConfidence: "multi-step-safe-auto-fix", AutoFixPolicy: "auto-fix"}
	case "single_quotes", "unquoted_keys", "missing_comma", "missing_colon", "missing_closing_delimiter", "unterminated_string":
		return policy{PrimaryRule: pattern, Severity: "error", FixConfidence: "suggestion-or-future-fix", AutoFixPolicy: "lint-only"}
	case "multiple_top_level_values":
		return policy{PrimaryRule: pattern, Severity: "error", FixConfidence: "suggestion-only", AutoFixPolicy: "lint-only"}
	case "duplicate_key":
		return policy{PrimaryRule: pattern, Severity: "warn", FixConfidence: "lint-only", AutoFixPolicy: "lint-only"}
	default:
		return policy{PrimaryRule: pattern, Severity: "info", FixConfidence: "n/a", AutoFixPolicy: "none"}
	}
}

func parserSignal(result jsonsanitize.Result) string {
	switch {
	case len(result.OriginalErrors) > 0 && !result.OriginalStrictParseClean:
		return "tree-sitter+strict"
	case len(result.OriginalErrors) > 0:
		return "tree-sitter"
	case !result.OriginalStrictParseClean:
		return "strict-only"
	default:
		return "clean"
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
	for name := range set {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func uniqueFixNames(fixes []jsonsanitize.Fix) []string {
	set := map[string]struct{}{}
	for _, fix := range fixes {
		set[fix.Rule] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for name := range set {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
