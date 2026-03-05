// Package sanitize provides YAML error detection and heuristic fixing using tree-sitter.
//
// Strategy:
//   - Parse the input with tree-sitter (tree-sitter-grammars/tree-sitter-yaml v0.7.2,
//     which is the actively maintained fork of ikatyang/tree-sitter-yaml with improved
//     error recovery, ABI 14, and scanner rewritten in C).
//   - Collect all ERROR and MISSING nodes (structural parse failures).
//   - Additionally run a set of "lint" rules that detect common YAML mistakes
//     even when tree-sitter does not raise an ERROR (e.g. duplicate keys,
//     missing space after colon, bad list dash, trailing comma, tab indent).
//   - Apply fixes iteratively until the tree is clean or no more progress.
package sanitize

import (
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	yaml "github.com/tree-sitter-grammars/tree-sitter-yaml/bindings/go"
)

// ErrorNode describes a tree-sitter ERROR or MISSING node found during parsing.
type ErrorNode struct {
	Type      string `json:"type"`       // "ERROR" or "MISSING"
	StartByte uint   `json:"start_byte"`
	EndByte   uint   `json:"end_byte"`
	StartRow  uint   `json:"start_row"`
	StartCol  uint   `json:"start_col"`
	EndRow    uint   `json:"end_row"`
	EndCol    uint   `json:"end_col"`
	Text      string `json:"text"`
}

// LintIssue describes a structural issue found by the linter (not necessarily a
// tree-sitter ERROR — e.g. duplicate keys, missing space after colon).
type LintIssue struct {
	Rule        string `json:"rule"`
	Description string `json:"description"`
	Row         int    `json:"row"` // 0-indexed
}

// Fix describes one applied heuristic fix.
type Fix struct {
	Rule        string `json:"rule"`
	Description string `json:"description"`
	Before      string `json:"before"`
	After       string `json:"after"`
}

// Result is the full output of a sanitize pass.
type Result struct {
	Original          string      `json:"original"`
	Sanitized         string      `json:"sanitized"`
	TreeText          string      `json:"tree_text"`           // tree of the sanitized output
	OriginalTreeText  string      `json:"original_tree_text"`  // tree of the original input
	Errors            []ErrorNode `json:"errors"`              // errors in sanitized output
	OriginalErrors    []ErrorNode `json:"original_errors"`     // errors in original input
	LintIssues        []LintIssue `json:"lint_issues"`
	OriginalLintIssues []LintIssue `json:"original_lint_issues"`
	Fixes             []Fix       `json:"fixes"`
	ParseClean        bool        `json:"parse_clean"`  // true if sanitized has no tree-sitter errors
	LintClean         bool        `json:"lint_clean"`   // true if no lint issues remain
}

// newParser creates a configured tree-sitter parser for YAML.
func newParser() *sitter.Parser {
	parser := sitter.NewParser()
	lang := sitter.NewLanguage(yaml.Language())
	if err := parser.SetLanguage(lang); err != nil {
		panic("failed to set YAML language: " + err.Error())
	}
	return parser
}

// ParseTree returns the tree-sitter sexp representation of the YAML source,
// plus any ERROR/MISSING nodes found.
func ParseTree(src string) (string, []ErrorNode, error) {
	parser := newParser()
	defer parser.Close()

	content := []byte(src)
	tree := parser.Parse(content, nil)
	if tree == nil {
		return "", nil, fmt.Errorf("tree-sitter returned nil tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	errors := collectErrors(root, content)
	return root.ToSexp(), errors, nil
}

// collectErrors walks the tree and gathers all ERROR and MISSING nodes.
func collectErrors(node *sitter.Node, src []byte) []ErrorNode {
	var errs []ErrorNode
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.IsError() || n.IsMissing() {
			typ := "ERROR"
			if n.IsMissing() {
				typ = "MISSING"
			}
			sp := n.StartPosition()
			ep := n.EndPosition()
			text := ""
			sb := n.StartByte()
			eb := n.EndByte()
			if sb < eb && int(eb) <= len(src) {
				text = string(src[sb:eb])
			}
			errs = append(errs, ErrorNode{
				Type:      typ,
				StartByte: sb,
				EndByte:   eb,
				StartRow:  sp.Row,
				StartCol:  sp.Column,
				EndRow:    ep.Row,
				EndCol:    ep.Column,
				Text:      text,
			})
		}
		count := n.ChildCount()
		cursor := n.Walk()
		defer cursor.Close()
		children := n.Children(cursor)
		for i := range children {
			walk(&children[i])
		}
		_ = count
	}
	walk(node)
	return errs
}

// ── Linter ────────────────────────────────────────────────────────────────

var (
	reTabLead       = regexp.MustCompile(`^\t+`)
	reMissingSpace  = regexp.MustCompile(`^(\s*)([^#\s][^:]*):([^\s/\n][^\n]*)$`)
	reBadDash       = regexp.MustCompile(`^(\s*)-([^\s\-\n])`)
	reTrailingComma = regexp.MustCompile(`,\s*([}\]])`)
	reKeyLine       = regexp.MustCompile(`^(\s*)([^#\s\-\[{][^:]*[^:\s]):\s`)
	reExtraColon    = regexp.MustCompile(`:\s`)
)

// Lint scans the source for common YAML mistakes and returns a list of issues.
func Lint(src string) []LintIssue {
	var issues []LintIssue
	lines := strings.Split(src, "\n")

	seenKeys := map[string]bool{}

	for i, line := range lines {
		row := i + 1 // 1-indexed for display

		// Tab indentation
		if reTabLead.MatchString(line) {
			issues = append(issues, LintIssue{
				Rule:        "tab_indent",
				Description: fmt.Sprintf("Line %d: tab used for indentation (YAML requires spaces)", row),
				Row:         i,
			})
		}

		// Missing space after colon
		if reMissingSpace.MatchString(line) {
			m := reMissingSpace.FindStringSubmatch(line)
			if m != nil {
				val := m[3]
				// Skip URLs (http://, https://)
				if !strings.HasPrefix(strings.TrimSpace(val), "//") {
					issues = append(issues, LintIssue{
						Rule:        "missing_space_after_colon",
						Description: fmt.Sprintf("Line %d: missing space after colon", row),
						Row:         i,
					})
				}
			}
		}

		// Bad list dash
		if reBadDash.MatchString(line) {
			issues = append(issues, LintIssue{
				Rule:        "list_dash_no_space",
				Description: fmt.Sprintf("Line %d: list dash not followed by a space", row),
				Row:         i,
			})
		}

		// Trailing comma in flow collection
		if reTrailingComma.MatchString(line) {
			issues = append(issues, LintIssue{
				Rule:        "trailing_comma",
				Description: fmt.Sprintf("Line %d: trailing comma in flow collection", row),
				Row:         i,
			})
		}

		// Duplicate keys (simple heuristic: track indent+key pairs)
		if m := reKeyLine.FindStringSubmatch(line); m != nil {
			indent := len(m[1])
			key := strings.TrimSpace(m[2])
			mapKey := fmt.Sprintf("%d:%s", indent, key)
			if seenKeys[mapKey] {
				issues = append(issues, LintIssue{
					Rule:        "duplicate_key",
					Description: fmt.Sprintf("Line %d: duplicate key '%s'", row, key),
					Row:         i,
				})
			}
			seenKeys[mapKey] = true
		}

		// Extra colon in plain scalar value
		// Pattern: key: some: value  (value contains ": ")
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "-") {
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
					reExtraColon.MatchString(rest) {
					issues = append(issues, LintIssue{
						Rule:        "extra_colon_in_value",
						Description: fmt.Sprintf("Line %d: plain scalar value contains a colon — may need quoting", row),
						Row:         i,
					})
				}
			}
		}
	}
	return issues
}

// ── Sanitizer ─────────────────────────────────────────────────────────────

// Sanitize attempts to fix common YAML errors heuristically, iterating until
// the tree is clean or no more fixes can be applied.
func Sanitize(src string) Result {
	original := src
	var allFixes []Fix

	// Capture original parse state before any fixes.
	origTreeText, origErrors, _ := ParseTree(original)
	origLintIssues := Lint(original)

	// Run up to 10 fix iterations.
	for iter := 0; iter < 10; iter++ {
		treeText, errors, err := ParseTree(src)
		if err != nil {
			break
		}
		lintIssues := Lint(src)

		if len(errors) == 0 && len(lintIssues) == 0 {
			return Result{
				Original:           original,
				Sanitized:          src,
				TreeText:           treeText,
				OriginalTreeText:   origTreeText,
				Errors:             errors,
				OriginalErrors:     origErrors,
				LintIssues:         lintIssues,
				OriginalLintIssues: origLintIssues,
				Fixes:              allFixes,
				ParseClean:         true,
				LintClean:          true,
			}
		}

		fixed, fixes := applyFixes(src, errors, lintIssues)
		if len(fixes) == 0 {
			// No progress — stop.
			treeText2, errors2, _ := ParseTree(src)
			lintIssues2 := Lint(src)
			return Result{
				Original:           original,
				Sanitized:          src,
				TreeText:           treeText2,
				OriginalTreeText:   origTreeText,
				Errors:             errors2,
				OriginalErrors:     origErrors,
				LintIssues:         lintIssues2,
				OriginalLintIssues: origLintIssues,
				Fixes:              allFixes,
				ParseClean:         len(errors2) == 0,
				LintClean:          len(lintIssues2) == 0,
			}
		}
		allFixes = append(allFixes, fixes...)
		src = fixed
	}

	treeText, errors, _ := ParseTree(src)
	lintIssues := Lint(src)
	return Result{
		Original:           original,
		Sanitized:          src,
		TreeText:           treeText,
		OriginalTreeText:   origTreeText,
		Errors:             errors,
		OriginalErrors:     origErrors,
		LintIssues:         lintIssues,
		OriginalLintIssues: origLintIssues,
		Fixes:              allFixes,
		ParseClean:         len(errors) == 0,
		LintClean:          len(lintIssues) == 0,
	}
}

// applyFixes applies one round of heuristic fixes to the source.
func applyFixes(src string, errors []ErrorNode, lintIssues []LintIssue) (string, []Fix) {
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

		newLine, f := fixLine(line, row, rules, hasTreeErr)
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
func fixLine(line string, row int, rules []string, hasTreeErr bool) (string, []Fix) {
	var fixes []Fix

	hasRule := func(r string) bool {
		for _, x := range rules {
			if x == r {
				return true
			}
		}
		return false
	}

	// Rule: tab indentation → 2 spaces per tab.
	if hasRule("tab_indent") {
		newLine := replaceLeadingTabs(line)
		if newLine != line {
			fixes = append(fixes, Fix{
				Rule:        "tab_indent",
				Description: fmt.Sprintf("Line %d: replaced leading tab(s) with 2 spaces each", row+1),
				Before:      line,
				After:       newLine,
			})
			line = newLine
		}
	}

	// Rule: missing space after colon  (a:1 → a: 1).
	if hasRule("missing_space_after_colon") {
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
	if hasRule("list_dash_no_space") {
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
	if hasRule("trailing_comma") {
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
	// Only apply when tree-sitter also sees an error OR lint flagged it.
	if hasRule("extra_colon_in_value") || hasTreeErr {
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
	re := regexp.MustCompile(`^(\s*[^:]+):([^\s/])`)
	return re.ReplaceAllStringFunc(line, func(m string) string {
		parts := re.FindStringSubmatch(m)
		if len(parts) < 3 {
			return m
		}
		return parts[1] + ": " + parts[2]
	})
}

// fixExtraColonInValue quotes the value part of a mapping line if it contains ": ".
func fixExtraColonInValue(line string) string {
	re := regexp.MustCompile(`^(\s*[^#\s][^:]*:\s+)([^'"{\[|>\n][^\n]*)$`)
	m := re.FindStringSubmatch(line)
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
	reKey := regexp.MustCompile(`^(\s*)([^#\s\-\[{][^:]*):\s`)
	for i, line := range lines {
		m := reKey.FindStringSubmatch(line)
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
			// Round DOWN to the nearest lower multiple of unit
			// (conservative: prefer same-level sibling over deeper child).
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

// replaceLeadingTabs replaces each leading tab with 2 spaces.
func replaceLeadingTabs(line string) string {
	i := 0
	for i < len(line) && line[i] == '\t' {
		i++
	}
	return strings.Repeat("  ", i) + line[i:]
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
