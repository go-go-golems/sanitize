package jsonsanitize

import (
	stdjson "encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func applyFixes(src string, doc documentAnalysis, cfg *config) (string, []Fix) {
	current := src
	var fixes []Fix

	if cfg.ruleEnabled("markdown_fence_wrapper") {
		next, fix, changed := fixMarkdownFenceWrapper(current)
		if changed {
			current = next
			fixes = append(fixes, fix)
		}
	}

	if cfg.ruleEnabled("leading_or_trailing_prose") {
		next, fix, changed := fixLeadingTrailingProse(current)
		if changed {
			current = next
			fixes = append(fixes, fix)
		}
	}

	if cfg.ruleEnabled("single_quotes") {
		next, roundFixes := fixSingleQuotedStrings(current)
		if next != current {
			current = next
			fixes = append(fixes, roundFixes...)
		}
	}

	if cfg.ruleEnabled("python_literals") {
		next, roundFixes := fixPythonLiterals(current)
		if next != current {
			current = next
			fixes = append(fixes, roundFixes...)
		}
	}

	signals := scanHeuristicSignals(current)
	if cfg.ruleEnabled("comment") {
		next, roundFixes := removeComments(current, signals.commentSpans)
		if next != current {
			current = next
			fixes = append(fixes, roundFixes...)
		}
	}

	signals = scanHeuristicSignals(current)
	if cfg.ruleEnabled("duplicate_comma") {
		next, roundFixes := removeDuplicateCommas(current, signals.duplicateCommaSpans)
		if next != current {
			current = next
			fixes = append(fixes, roundFixes...)
		}
	}

	signals = scanHeuristicSignals(current)
	if cfg.ruleEnabled("trailing_comma") {
		next, roundFixes := removeTrailingCommas(current, signals.trailingCommaSpans)
		if next != current {
			current = next
			fixes = append(fixes, roundFixes...)
		}
	}

	return current, fixes
}

func fixMarkdownFenceWrapper(src string) (string, Fix, bool) {
	if !reMarkdownFence.MatchString(src) {
		return src, Fix{}, false
	}

	lines := strings.Split(src, "\n")
	if len(lines) < 2 {
		return src, Fix{}, false
	}

	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	if start >= len(lines) || !strings.HasPrefix(strings.TrimSpace(lines[start]), "```") {
		return src, Fix{}, false
	}

	end := len(lines) - 1
	for end >= 0 && strings.TrimSpace(lines[end]) == "" {
		end--
	}
	if end <= start || strings.TrimSpace(lines[end]) != "```" {
		return src, Fix{}, false
	}

	body := strings.Join(lines[start+1:end], "\n")
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}

	return body, Fix{
		Rule:        "markdown_fence_wrapper",
		Description: "Removed Markdown code fence wrapper around JSON payload",
		Before:      src,
		After:       body,
	}, true
}

func fixLeadingTrailingProse(src string) (string, Fix, bool) {
	trimmedStart := firstNonWhitespaceIndex(src)
	trimmedEnd := lastNonWhitespaceIndex(src)
	if trimmedStart == -1 || trimmedEnd == -1 {
		return src, Fix{}, false
	}

	firstStructural := strings.IndexAny(src, "{[")
	lastStructural := strings.LastIndexAny(src, "}]")
	if firstStructural == -1 || lastStructural == -1 || firstStructural >= lastStructural {
		return src, Fix{}, false
	}
	if firstStructural == trimmedStart && lastStructural == trimmedEnd {
		return src, Fix{}, false
	}

	extracted := strings.TrimSpace(src[firstStructural : lastStructural+1])
	if extracted == "" {
		return src, Fix{}, false
	}
	extracted += "\n"

	return extracted, Fix{
		Rule:        "leading_or_trailing_prose",
		Description: "Extracted the likely JSON payload from surrounding prose",
		Before:      src,
		After:       extracted,
	}, true
}

func fixPythonLiterals(src string) (string, []Fix) {
	var out strings.Builder
	out.Grow(len(src))

	var fixes []Fix
	var quote byte
	escape := false

	for i := 0; i < len(src); {
		ch := src[i]

		if quote != 0 {
			out.WriteByte(ch)
			if escape {
				escape = false
			} else if ch == '\\' {
				escape = true
			} else if ch == quote {
				quote = 0
			}
			i++
			continue
		}

		if ch == '"' || ch == '\'' {
			quote = ch
			out.WriteByte(ch)
			i++
			continue
		}

		replaced, replacement, ok := matchPythonLiteral(src[i:])
		if ok && literalBoundaryBefore(src, i) && literalBoundaryAfter(src, i+len(replaced)) {
			out.WriteString(replacement)
			fixes = append(fixes, Fix{
				Rule:        "python_literals",
				Description: fmt.Sprintf("Normalized Python literal %q to %q", replaced, replacement),
				Before:      replaced,
				After:       replacement,
			})
			i += len(replaced)
			continue
		}

		out.WriteByte(ch)
		i++
	}

	return out.String(), fixes
}

func fixSingleQuotedStrings(src string) (string, []Fix) {
	var out strings.Builder
	out.Grow(len(src))

	var fixes []Fix
	var quote byte
	escape := false

	for i := 0; i < len(src); {
		ch := src[i]

		if quote == '"' {
			out.WriteByte(ch)
			if escape {
				escape = false
			} else if ch == '\\' {
				escape = true
			} else if ch == '"' {
				quote = 0
			}
			i++
			continue
		}

		if ch == '"' {
			quote = '"'
			out.WriteByte(ch)
			i++
			continue
		}

		if ch == '\'' && singleQuotedContextBefore(src, i) {
			raw, decoded, end, ok := parseSingleQuotedString(src, i)
			if ok && singleQuotedContextAfter(src, end) {
				replacementBytes, err := stdjson.Marshal(decoded)
				if err == nil {
					replacement := string(replacementBytes)
					out.WriteString(replacement)
					fixes = append(fixes, Fix{
						Rule:        "single_quotes",
						Description: "Rewrote single-quoted JSON string using double quotes",
						Before:      raw,
						After:       replacement,
					})
					i = end
					continue
				}
			}
		}

		out.WriteByte(ch)
		i++
	}

	return out.String(), fixes
}

func matchPythonLiteral(src string) (string, string, bool) {
	switch {
	case strings.HasPrefix(src, "True"):
		return "True", "true", true
	case strings.HasPrefix(src, "False"):
		return "False", "false", true
	case strings.HasPrefix(src, "None"):
		return "None", "null", true
	default:
		return "", "", false
	}
}

func literalBoundaryBefore(src string, idx int) bool {
	if idx <= 0 {
		return true
	}
	r := rune(src[idx-1])
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
}

func literalBoundaryAfter(src string, idx int) bool {
	if idx >= len(src) {
		return true
	}
	r := rune(src[idx])
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
}

func parseSingleQuotedString(src string, start int) (string, string, int, bool) {
	if start < 0 || start >= len(src) || src[start] != '\'' {
		return "", "", 0, false
	}

	var out strings.Builder
	for i := start + 1; i < len(src); {
		ch := src[i]
		if ch == '\'' {
			return src[start : i+1], out.String(), i + 1, true
		}
		if ch != '\\' {
			out.WriteByte(ch)
			i++
			continue
		}

		if i+1 >= len(src) {
			return "", "", 0, false
		}

		switch esc := src[i+1]; esc {
		case '\'', '"', '\\', '/':
			out.WriteByte(esc)
			i += 2
		case 'b':
			out.WriteByte('\b')
			i += 2
		case 'f':
			out.WriteByte('\f')
			i += 2
		case 'n':
			out.WriteByte('\n')
			i += 2
		case 'r':
			out.WriteByte('\r')
			i += 2
		case 't':
			out.WriteByte('\t')
			i += 2
		case 'u':
			if i+6 > len(src) {
				return "", "", 0, false
			}
			value, err := strconv.ParseUint(src[i+2:i+6], 16, 64)
			if err != nil {
				return "", "", 0, false
			}
			out.WriteRune(rune(value))
			i += 6
		default:
			out.WriteByte('\\')
			out.WriteByte(esc)
			i += 2
		}
	}

	return "", "", 0, false
}

func singleQuotedContextBefore(src string, idx int) bool {
	prev := previousNonSpaceByte(src, idx-1)
	switch prev {
	case 0, '{', '[', ':', ',':
		return true
	default:
		return false
	}
}

func singleQuotedContextAfter(src string, idx int) bool {
	next := nextNonSpaceByte(src, idx)
	switch next {
	case 0, ':', ',', '}', ']':
		return true
	default:
		return false
	}
}

func previousNonSpaceByte(src string, idx int) byte {
	for i := idx; i >= 0; i-- {
		if !isJSONSpace(src[i]) {
			return src[i]
		}
	}
	return 0
}

func nextNonSpaceByte(src string, idx int) byte {
	for i := idx; i < len(src); i++ {
		if !isJSONSpace(src[i]) {
			return src[i]
		}
	}
	return 0
}

func isJSONSpace(ch byte) bool {
	switch ch {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

func removeComments(src string, spans []span) (string, []Fix) {
	return rewriteSpans(src, spans, "comment", "Removed comment syntax from JSON input", func(_ string) string {
		return ""
	})
}

func removeDuplicateCommas(src string, spans []span) (string, []Fix) {
	return rewriteSpans(src, spans, "duplicate_comma", "Collapsed duplicate comma in JSON input", func(_ string) string {
		return ","
	})
}

func removeTrailingCommas(src string, spans []span) (string, []Fix) {
	return rewriteSpans(src, spans, "trailing_comma", "Removed trailing comma before closing delimiter", func(segment string) string {
		if segment == "" {
			return segment
		}
		return segment[len(segment)-1:]
	})
}

func rewriteSpans(src string, spans []span, rule, description string, replacer func(string) string) (string, []Fix) {
	if len(spans) == 0 {
		return src, nil
	}

	sort.Slice(spans, func(i, j int) bool {
		return spans[i].start < spans[j].start
	})

	var out strings.Builder
	out.Grow(len(src))

	var fixes []Fix
	last := 0
	for _, sp := range spans {
		if sp.start < last || sp.end < sp.start || sp.end > len(src) {
			continue
		}
		out.WriteString(src[last:sp.start])
		before := src[sp.start:sp.end]
		after := replacer(before)
		out.WriteString(after)
		fixes = append(fixes, Fix{
			Rule:        rule,
			Description: description,
			Before:      before,
			After:       after,
		})
		last = sp.end
	}
	out.WriteString(src[last:])
	return out.String(), fixes
}
