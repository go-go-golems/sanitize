package jsonsanitize

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
)

// Lint scans the source for JSON issues and returns a list of issues.
func Lint(src string) []LintIssue {
	issues, err := LintWithOptions(src)
	if err != nil {
		return nil
	}
	return issues
}

// LintWithOptions scans the source for JSON issues using the provided rule configuration.
func LintWithOptions(src string, opts ...Option) ([]LintIssue, error) {
	cfg, err := buildConfig(opts...)
	if err != nil {
		return nil, err
	}

	doc, err := analyzeDocument(src)
	if err != nil {
		return nil, err
	}
	return lintIssuesFromAnalysis(doc, &cfg), nil
}

func lintIssuesFromAnalysis(doc documentAnalysis, cfg *config) []LintIssue {
	issues := lintFromParseErrors(doc.ParseErrors, cfg)
	issues = append(issues, lintFromStrictParseError(doc, cfg)...)
	issues = append(issues, lintFromDuplicateKeys(doc, cfg)...)
	return issues
}

func lintFromParseErrors(errors []ErrorNode, cfg *config) []LintIssue {
	issues := make([]LintIssue, 0, len(errors))
	for _, parseErr := range errors {
		rule := "structural_parse_error"
		description := fmt.Sprintf(
			"Line %d: JSON structure could not be parsed cleanly",
			parseErr.StartRow+1,
		)
		if parseErr.Type == "MISSING" {
			rule = "missing_syntax_node"
			description = fmt.Sprintf(
				"Line %d: JSON syntax appears to be missing a required element",
				parseErr.StartRow+1,
			)
		}
		if !cfg.ruleEnabled(rule) {
			continue
		}

		issues = append(issues, LintIssue{
			Rule:        rule,
			Source:      "parse",
			Description: description,
			StartByte:   parseErr.StartByte,
			EndByte:     parseErr.EndByte,
			StartRow:    parseErr.StartRow,
			StartCol:    parseErr.StartCol,
			EndRow:      parseErr.EndRow,
			EndCol:      parseErr.EndCol,
			Row:         int(parseErr.StartRow),
		})
	}
	return issues
}

func lintFromStrictParseError(doc documentAnalysis, cfg *config) []LintIssue {
	if doc.StrictParseError == nil || !cfg.ruleEnabled("strict_parse_error") {
		return nil
	}

	startByte, endByte := strictParseIssueSpan(doc)
	startRow, startCol := doc.LineIndex.rowColAtByte(startByte)
	endRow, endCol := doc.LineIndex.rowColAtByte(endByte)

	return []LintIssue{{
		Rule:        "strict_parse_error",
		Source:      "strict-parser",
		Description: fmt.Sprintf("Line %d: strict JSON parse failed: %s", startRow+1, doc.StrictParseError.Error()),
		StartByte:   startByte,
		EndByte:     endByte,
		StartRow:    uint(startRow),
		StartCol:    uint(startCol),
		EndRow:      uint(endRow),
		EndCol:      uint(endCol),
		Row:         startRow,
	}}
}

func lintFromDuplicateKeys(doc documentAnalysis, cfg *config) []LintIssue {
	if !cfg.ruleEnabled("duplicate_key") {
		return nil
	}

	issues := make([]LintIssue, 0, len(doc.DuplicateKeys))
	for _, duplicate := range doc.DuplicateKeys {
		startRow, startCol := doc.LineIndex.rowColAtByte(duplicate.StartByte)
		endRow, endCol := doc.LineIndex.rowColAtByte(duplicate.EndByte)
		issues = append(issues, LintIssue{
			Rule:        "duplicate_key",
			Source:      "heuristic",
			Description: fmt.Sprintf("Line %d: duplicate key '%s'", startRow+1, duplicate.Key),
			StartByte:   duplicate.StartByte,
			EndByte:     duplicate.EndByte,
			StartRow:    uint(startRow),
			StartCol:    uint(startCol),
			EndRow:      uint(endRow),
			EndCol:      uint(endCol),
			Row:         startRow,
		})
	}
	return issues
}

func strictParseIssueSpan(doc documentAnalysis) (uint, uint) {
	var syntaxErr *stdjson.SyntaxError
	if errors.As(doc.StrictParseError, &syntaxErr) && syntaxErr.Offset > 0 {
		start := uint(syntaxErr.Offset - 1)
		return start, clampByte(start+1, doc.LineIndex)
	}

	var multiErr multiValueError
	if errors.As(doc.StrictParseError, &multiErr) && multiErr.Offset > 0 {
		start := uint(multiErr.Offset - 1)
		return start, clampByte(start+1, doc.LineIndex)
	}

	return 0, 0
}

func clampByte(offset uint, li lineIndex) uint {
	if len(li.starts) == 0 {
		return offset
	}
	lastStart := li.starts[len(li.starts)-1]
	if int(offset) < lastStart {
		return offset
	}
	return offset
}
