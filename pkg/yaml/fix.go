package yamlsanitize

import (
	"fmt"
	"regexp"
	"strings"
)

// Precompiled regexes for fixers.
var (
	reMissingSpaceFix = regexp.MustCompile(`^(\s*[^:]+):([^\s/])`)
	reExtraColonInVal = regexp.MustCompile(`^(\s*[^#\s][^:]*:\s+)([^'"{\[|>\n][^\n]*)$`)
)

// applyFixes applies one round of heuristic fixes to the source.
func applyFixes(src string, doc documentAnalysis, cfg *config) (string, []Fix) {
	var fixes []Fix
	lines := strings.Split(src, "\n")
	lintIssues := lintIssuesFromAnalysis(src, doc, cfg)

	// Build a set of lint rules per row for quick lookup.
	lintByRow := map[int][]string{}
	for _, li := range lintIssues {
		lintByRow[li.Row] = append(lintByRow[li.Row], li.Rule)
	}

	// Build error rows set from full parse-error spans.
	errorRows := map[int]bool{}
	for _, e := range doc.ParseErrors {
		for row := int(e.StartRow); row <= int(e.EndRow); row++ {
			errorRows[row] = true
		}
	}

	changed := false
	for row := range lines {
		line := lines[row]
		rules := lintByRow[row]
		hasTreeErr := errorRows[row]

		newLine, f := fixLine(line, row, rules, hasTreeErr, cfg)
		if newLine != line {
			lines[row] = newLine
			fixes = append(fixes, f...)
			changed = true
		}
	}

	// Document-level fixes (duplicate keys).
	if !changed && cfg.ruleEnabled("duplicate_key") {
		newSrc, f := fixDuplicateKeysOccurrences(strings.Join(lines, "\n"), doc.DuplicateKeys)
		if newSrc != strings.Join(lines, "\n") {
			return newSrc, f
		}
	}

	// Document-level fix: mixed/inconsistent indentation depth.
	// Only attempt when tree-sitter reported errors (structural breakage).
	if !changed && cfg.ruleEnabled("mixed_indent") && len(doc.ParseErrors) > 0 {
		newSrc, f := fixMixedIndentation(strings.Join(lines, "\n"))
		if len(f) > 0 {
			return newSrc, f
		}
	}

	return strings.Join(lines, "\n"), fixes
}

// fixLine applies heuristic fixes to a single line.
func fixLine(line string, row int, rules []string, hasTreeErr bool, cfg *config) (string, []Fix) {
	var fixes []Fix

	hasRule := func(r string) bool {
		for _, x := range rules {
			if x == r {
				return true
			}
		}
		return false
	}

	// Rule: tab indentation → spaces per tab.
	if hasRule("tab_indent") && cfg.ruleEnabled("tab_indent") {
		newLine := replaceLeadingTabs(line, cfg.tabWidth)
		if newLine != line {
			fixes = append(fixes, Fix{
				Rule:        "tab_indent",
				Description: fmt.Sprintf("Line %d: replaced leading tab(s) with %d spaces each", row+1, cfg.tabWidth),
				Before:      line,
				After:       newLine,
			})
			line = newLine
		}
	}

	// Rule: missing space after colon  (a:1 → a: 1).
	if hasRule("missing_space_after_colon") && cfg.ruleEnabled("missing_space_after_colon") {
		newLine := fixMissingSpaceAfterColon(line)
		if newLine != line {
			fixes = append(fixes, Fix{
				Rule:        "missing_space_after_colon",
				Description: fmt.Sprintf("Line %d: added space after colon", row+1),
				Before:      line,
				After:       newLine,
			})
			line = newLine
		}
	}

	// Rule: list dash not followed by space  (-a → - a).
	if hasRule("list_dash_no_space") && cfg.ruleEnabled("list_dash_no_space") {
		newLine := reBadDash.ReplaceAllString(line, "${1}- ${2}")
		if newLine != line {
			fixes = append(fixes, Fix{
				Rule:        "list_dash_no_space",
				Description: fmt.Sprintf("Line %d: added space after list dash", row+1),
				Before:      line,
				After:       newLine,
			})
			line = newLine
		}
	}

	// Rule: trailing comma in flow collection.
	if hasRule("trailing_comma") && cfg.ruleEnabled("trailing_comma") {
		newLine := reTrailingComma.ReplaceAllString(line, "$1")
		if newLine != line {
			fixes = append(fixes, Fix{
				Rule:        "trailing_comma",
				Description: fmt.Sprintf("Line %d: removed trailing comma in flow collection", row+1),
				Before:      line,
				After:       newLine,
			})
			line = newLine
		}
	}

	// Rule: extra colon in plain scalar value — quote the value.
	if hasTreeErr && cfg.ruleEnabled("extra_colon_in_value") {
		newLine := fixExtraColonInValue(line)
		if newLine != line {
			fixes = append(fixes, Fix{
				Rule:        "extra_colon_in_value",
				Description: fmt.Sprintf("Line %d: quoted value containing extra colon", row+1),
				Before:      line,
				After:       newLine,
			})
			line = newLine
		}
	}

	return line, fixes
}

// fixMissingSpaceAfterColon adds a space after the first colon in a mapping line.
func fixMissingSpaceAfterColon(line string) string {
	return reMissingSpaceFix.ReplaceAllStringFunc(line, func(m string) string {
		parts := reMissingSpaceFix.FindStringSubmatch(m)
		if len(parts) < 3 {
			return m
		}
		return parts[1] + ": " + parts[2]
	})
}

// fixExtraColonInValue quotes the value part of a mapping line if it contains ": ".
func fixExtraColonInValue(line string) string {
	m := reExtraColonInVal.FindStringSubmatch(line)
	if m == nil {
		return line
	}
	val := strings.TrimRight(m[2], " \t")
	// Only quote if value contains ": " or ends with ":"
	if !reExtraColon.MatchString(val) && !strings.HasSuffix(val, ":") {
		return line
	}
	// Don't double-quote
	if strings.HasPrefix(val, `"`) || strings.HasPrefix(val, `'`) {
		return line
	}
	return m[1] + quoteValue(val)
}

func fixDuplicateKeysOccurrences(src string, duplicates []duplicateKeyOccurrence) (string, []Fix) {
	if len(duplicates) == 0 {
		return src, nil
	}

	lines := strings.Split(src, "\n")
	fixes := make([]Fix, 0, len(duplicates))

	for i := len(duplicates) - 1; i >= 0; i-- {
		duplicate := duplicates[i]
		if duplicate.EndByte > uint(len(src)) || duplicate.StartByte >= duplicate.EndByte {
			continue
		}

		newKeyText := duplicateKeyReplacement(duplicate.KeyText, duplicate.DuplicateIndex)
		src = src[:duplicate.StartByte] + newKeyText + src[duplicate.EndByte:]

		line := lineIndexAtByte(src, duplicate.StartByte)
		if line >= 0 && line < len(lines) {
			before := lines[line]
			after := strings.Replace(before, duplicate.KeyText, newKeyText, 1)
			lines[line] = after
			fixes = append([]Fix{{
				Rule:        "duplicate_key",
				Description: fmt.Sprintf("Line %d: renamed duplicate key '%s' → '%s'", line+1, duplicate.Key, duplicateKeyIdentity(newKeyText)),
				Before:      before,
				After:       after,
			}}, fixes...)
		}
	}

	return src, fixes
}

// fixMixedIndentation detects the dominant indent width and normalises lines
// whose leading-space count is not a multiple of that width.
func fixMixedIndentation(src string) (string, []Fix) {
	lines := strings.Split(src, "\n")
	unit, offenders := detectMixedIndentationRows(lines)
	if len(offenders) == 0 {
		return src, nil
	}

	// Re-indent lines whose space count is not a multiple of unit.
	var fixes []Fix
	newLines := make([]string, len(lines))
	copy(newLines, lines)
	offenderSet := map[int]bool{}
	for _, row := range offenders {
		offenderSet[row] = true
	}
	for i, line := range lines {
		if !offenderSet[i] {
			continue
		}

		spaces := leadingSpaces(line)
		normSpaces := (spaces / unit) * unit
		if normSpaces == 0 {
			normSpaces = unit
		}
		newLine := strings.Repeat(" ", normSpaces) + line[spaces:]
		newLines[i] = newLine
		fixes = append(fixes, Fix{
			Rule:        "mixed_indent",
			Description: fmt.Sprintf("Line %d: normalised indent %d→%d spaces (unit=%d)", i+1, spaces, normSpaces, unit),
			Before:      line,
			After:       newLine,
		})
	}
	return strings.Join(newLines, "\n"), fixes
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// replaceLeadingTabs replaces each leading tab with tabWidth spaces.
func replaceLeadingTabs(line string, tabWidth int) string {
	i := 0
	for i < len(line) && line[i] == '\t' {
		i++
	}
	return strings.Repeat(" ", i*tabWidth) + line[i:]
}

// quoteValue wraps a value in double quotes, escaping internal double quotes.
func quoteValue(v string) string {
	if (strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`)) ||
		(strings.HasPrefix(v, `'`) && strings.HasSuffix(v, `'`)) {
		return v
	}
	escaped := strings.ReplaceAll(v, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}
