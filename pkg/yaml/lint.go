package yamlsanitize

import (
	"fmt"
	"regexp"
	"strings"
)

// Precompiled regexes for lint rules.
var (
	reTabLead       = regexp.MustCompile(`^\t+`)
	reMissingSpace  = regexp.MustCompile(`^(\s*)([^#\s][^:]*):([^\s/\n][^\n]*)$`)
	reBadDash       = regexp.MustCompile(`^(\s*)-([^\s\-\n])`)
	reTrailingComma = regexp.MustCompile(`,\s*([}\]])`)
	reExtraColon    = regexp.MustCompile(`:\s`)
)

// Lint scans the source for common YAML mistakes and returns a list of issues.
func Lint(src string) []LintIssue {
	issues, err := LintWithOptions(src)
	if err != nil {
		return nil
	}
	return issues
}

// LintWithOptions scans the source for YAML issues using the provided rule configuration.
func LintWithOptions(src string, opts ...Option) ([]LintIssue, error) {
	cfg, err := buildConfig(opts...)
	if err != nil {
		return nil, err
	}
	return lintWithConfig(src, cfg), nil
}

func lintWithConfig(src string, cfg config) []LintIssue {
	doc, err := analyzeDocument(src)
	if err != nil {
		return lintLineIssues(src, lineIndex{}, &cfg, nil)
	}
	return lintIssuesFromAnalysis(src, doc, &cfg)
}

func lintIssuesFromAnalysis(src string, doc documentAnalysis, cfg *config) []LintIssue {
	issues := lintFromParseErrors(doc.ParseErrors, cfg)
	issues = append(issues, lintLineIssues(src, doc.LineIndex, cfg, doc.ParseErrors, doc.DuplicateKeys)...)
	issues = append(issues, lintMixedIndentation(src, doc, cfg)...)
	return issues
}

func lintLineIssues(src string, li lineIndex, cfg *config, parseErrors []ErrorNode, duplicates ...[]duplicateKeyOccurrence) []LintIssue {
	var issues []LintIssue
	lines := strings.Split(src, "\n")
	lineOffset := 0

	for i, line := range lines {
		row := i + 1 // 1-indexed for display
		startByte := uint(lineOffset)
		endByte := uint(lineOffset + len(line))

		// Tab indentation
		if cfg.ruleEnabled("tab_indent") && reTabLead.MatchString(line) {
			issues = append(issues, newLineLintIssue(
				"tab_indent",
				fmt.Sprintf("Line %d: tab used for indentation (YAML requires spaces)", row),
				i,
				startByte,
				endByte,
			))
		}

		// Missing space after colon
		if cfg.ruleEnabled("missing_space_after_colon") && reMissingSpace.MatchString(line) {
			m := reMissingSpace.FindStringSubmatch(line)
			if m != nil {
				val := m[3]
				// Skip URLs (http://, https://)
				if !strings.HasPrefix(strings.TrimSpace(val), "//") {
					issues = append(issues, newLineLintIssue(
						"missing_space_after_colon",
						fmt.Sprintf("Line %d: missing space after colon", row),
						i,
						startByte,
						endByte,
					))
				}
			}
		}

		// Bad list dash
		if cfg.ruleEnabled("list_dash_no_space") && reBadDash.MatchString(line) {
			issues = append(issues, newLineLintIssue(
				"list_dash_no_space",
				fmt.Sprintf("Line %d: list dash not followed by a space", row),
				i,
				startByte,
				endByte,
			))
		}

		// Trailing comma in flow collection
		if cfg.ruleEnabled("trailing_comma") && reTrailingComma.MatchString(line) {
			issues = append(issues, newLineLintIssue(
				"trailing_comma",
				fmt.Sprintf("Line %d: trailing comma in flow collection", row),
				i,
				startByte,
				endByte,
			))
		}

		// Extra colon in plain scalar value
		trimmed := strings.TrimSpace(line)
		if cfg.ruleEnabled("extra_colon_in_value") && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "-") {
			colonIdx := strings.Index(line, ": ")
			if colonIdx >= 0 {
				rest := strings.TrimSpace(line[colonIdx+2:])
				if len(rest) > 0 &&
					!strings.HasPrefix(rest, `"`) &&
					!strings.HasPrefix(rest, `'`) &&
					!strings.HasPrefix(rest, `{`) &&
					!strings.HasPrefix(rest, `[`) &&
					!strings.HasPrefix(rest, `|`) &&
					!strings.HasPrefix(rest, `>`) &&
					reExtraColon.MatchString(rest) &&
					parseErrorTouchesRow(parseErrors, i, 1) {
					issues = append(issues, newLineLintIssue(
						"extra_colon_in_value",
						fmt.Sprintf("Line %d: plain scalar value contains a colon — may need quoting", row),
						i,
						startByte,
						endByte,
					))
				}
			}
		}

		lineOffset += len(line)
		if i < len(lines)-1 {
			lineOffset++
		}
	}

	var duplicateKeys []duplicateKeyOccurrence
	if len(duplicates) > 0 {
		duplicateKeys = duplicates[0]
	}
	for _, duplicate := range duplicateKeys {
		if !cfg.ruleEnabled("duplicate_key") {
			continue
		}
		startRow, startCol := li.rowColAtByte(duplicate.StartByte)
		endRow, endCol := li.rowColAtByte(duplicate.EndByte)
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

func parseErrorTouchesRow(errors []ErrorNode, row int, nearbyDistance int) bool {
	for _, parseErr := range errors {
		startRow := int(parseErr.StartRow)
		endRow := int(parseErr.EndRow)
		if row >= startRow && row <= endRow {
			return true
		}
		if nearbyDistance > 0 && row >= startRow-nearbyDistance && row <= endRow+nearbyDistance {
			return true
		}
	}
	return false
}

func lintFromParseErrors(errors []ErrorNode, cfg *config) []LintIssue {
	issues := make([]LintIssue, 0, len(errors))
	for _, parseErr := range errors {
		rule := "structural_parse_error"
		description := fmt.Sprintf(
			"Line %d: YAML structure could not be parsed cleanly",
			parseErr.StartRow+1,
		)
		if parseErr.Type == "MISSING" {
			rule = "missing_syntax_node"
			description = fmt.Sprintf(
				"Line %d: YAML syntax appears to be missing a required element",
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

func newLineLintIssue(rule, description string, row int, startByte, endByte uint) LintIssue {
	return LintIssue{
		Rule:        rule,
		Source:      "heuristic",
		Description: description,
		StartByte:   startByte,
		EndByte:     endByte,
		StartRow:    uint(row),
		StartCol:    0,
		EndRow:      uint(row),
		EndCol:      endByte - startByte,
		Row:         row,
	}
}

func lintMixedIndentation(src string, doc documentAnalysis, cfg *config) []LintIssue {
	if !cfg.ruleEnabled("mixed_indent") {
		return nil
	}
	if len(doc.ParseErrors) == 0 {
		return nil
	}

	lines := strings.Split(src, "\n")
	unit, offenders := detectMixedIndentationRows(lines)
	if len(offenders) == 0 {
		return nil
	}

	issues := make([]LintIssue, 0, len(offenders))
	lineOffset := 0
	for i, line := range lines {
		startByte := uint(lineOffset)
		endByte := uint(lineOffset + len(line))

		for _, offender := range offenders {
			if offender == i {
				issues = append(issues, newLineLintIssue(
					"mixed_indent",
					fmt.Sprintf("Line %d: indentation is not a multiple of the dominant %d-space unit", i+1, unit),
					i,
					startByte,
					endByte,
				))
				break
			}
		}

		lineOffset += len(line)
		if i < len(lines)-1 {
			lineOffset++
		}
	}

	return issues
}
