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
