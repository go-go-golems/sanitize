package yamlsanitize

import (
	"fmt"
	"regexp"
	"strings"
)

// Precompiled regexes for fixers.
var (
	reMissingSpaceFix  = regexp.MustCompile(`^(\s*[^:]+):([^\s/])`)
	reExtraColonInVal  = regexp.MustCompile(`^(\s*[^#\s][^:]*:\s+)([^'"{\[|>\n][^\n]*)$`)
	reDuplicateKeyLine = regexp.MustCompile(`^(\s*)([^#\s\-\[{][^:]*):\s`)
)

// applyFixes applies one round of heuristic fixes to the source.
func applyFixes(src string, errors []ErrorNode, lintIssues []LintIssue, cfg *config) (string, []Fix) {
	var fixes []Fix
	lines := strings.Split(src, "\n")

	// Build a set of lint rules per row for quick lookup.
	lintByRow := map[int][]string{}
	for _, li := range lintIssues {
		lintByRow[li.Row] = append(lintByRow[li.Row], li.Rule)
	}

	// Build error rows set.
	errorRows := map[uint]bool{}
	for _, e := range errors {
		errorRows[e.StartRow] = true
	}

	changed := false
	for row := range lines {
		line := lines[row]
		rules := lintByRow[row]
		hasTreeErr := errorRows[uint(row)]

		newLine, f := fixLine(line, row, rules, hasTreeErr, cfg)
		if newLine != line {
			lines[row] = newLine
			fixes = append(fixes, f...)
			changed = true
		}
	}

	// Document-level fixes (duplicate keys).
	if !changed {
		newSrc, f := fixDuplicateKeys(strings.Join(lines, "\n"))
		if newSrc != strings.Join(lines, "\n") {
			return newSrc, f
		}
	}

	// Document-level fix: mixed/inconsistent indentation depth.
	// Only attempt when tree-sitter reported errors (structural breakage).
	if !changed && len(errors) > 0 {
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
	if (hasRule("extra_colon_in_value") || hasTreeErr) && cfg.ruleEnabled("extra_colon_in_value") {
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

// fixDuplicateKeys renames duplicate sibling keys by appending a numeric suffix.
func fixDuplicateKeys(src string) (string, []Fix) {
	var fixes []Fix
	lines := strings.Split(src, "\n")
	seen := map[string]int{}
	for i, line := range lines {
		m := reDuplicateKeyLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		indent := m[1]
		key := strings.TrimSpace(m[2])
		mapKey := indent + "|" + key
		if seen[mapKey] > 0 {
			newKey := fmt.Sprintf("%s_%d", key, seen[mapKey]+1)
			newLine := strings.Replace(line, m[2], newKey, 1)
			fixes = append(fixes, Fix{
				Rule:        "duplicate_key",
				Description: fmt.Sprintf("Line %d: renamed duplicate key '%s' → '%s'", i+1, key, newKey),
				Before:      line,
				After:       newLine,
			})
			lines[i] = newLine
		}
		seen[mapKey]++
	}
	return strings.Join(lines, "\n"), fixes
}

// fixMixedIndentation detects the dominant indent width and normalises lines
// whose leading-space count is not a multiple of that width.
func fixMixedIndentation(src string) (string, []Fix) {
	lines := strings.Split(src, "\n")

	// Count leading spaces per line (skip blank/comment lines).
	indentCounts := map[int]int{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		spaces := 0
		for _, ch := range line {
			if ch == ' ' {
				spaces++
			} else {
				break
			}
		}
		if spaces > 0 {
			indentCounts[spaces]++
		}
	}

	if len(indentCounts) == 0 {
		return src, nil
	}

	// Find the GCD of all observed indent widths as the dominant unit.
	gcdAll := 0
	for w := range indentCounts {
		gcdAll = gcd(gcdAll, w)
	}
	if gcdAll <= 0 {
		gcdAll = 2
	}
	// Prefer 2 or 4 if they divide gcdAll evenly.
	unit := gcdAll
	if unit == 1 {
		// Ambiguous — pick the most common indent width.
		best, bestCnt := 2, 0
		for w, cnt := range indentCounts {
			if cnt > bestCnt {
				best, bestCnt = w, cnt
			}
		}
		unit = best
	}

	// Re-indent lines whose space count is not a multiple of unit.
	var fixes []Fix
	newLines := make([]string, len(lines))
	copy(newLines, lines)
	for i, line := range lines {
		spaces := 0
		for _, ch := range line {
			if ch == ' ' {
				spaces++
			} else {
				break
			}
		}
		if spaces > 0 && spaces%unit != 0 {
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
