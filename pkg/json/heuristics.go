package jsonsanitize

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	reMarkdownFence = regexp.MustCompile("(?s)^\\s*```(?:json)?\\s*\\n.*\\n```\\s*$")
	reSingleKey     = regexp.MustCompile(`[{,]\s*'([^'\\]|\\.)*'\s*:`)
	reSingleValue   = regexp.MustCompile(`:\s*'([^'\\]|\\.)*'(\s*[,}\]])`)
	reUnquotedKey   = regexp.MustCompile(`[{,]\s*[A-Za-z_][A-Za-z0-9_-]*\s*:`)
	rePythonLiteral = regexp.MustCompile(`\b(True|False|None)\b`)
	rePlaceholder   = regexp.MustCompile(`<[^>\n]+>`)
)

type span struct {
	start int
	end   int
}

type opener struct {
	ch  byte
	pos int
}

type heuristicSignals struct {
	commentSpans        []span
	trailingCommaSpans  []span
	duplicateCommaSpans []span
	missingCloseOpener  *opener
	ellipsisSpans       []span
	placeholderSpans    []span
}

func lintFromHeuristics(src string, doc documentAnalysis, cfg *config) []LintIssue {
	signals := scanHeuristicSignals(src)
	var issues []LintIssue

	if cfg.ruleEnabled("markdown_fence_wrapper") && reMarkdownFence.MatchString(src) {
		issues = append(issues, newHeuristicIssue(
			"markdown_fence_wrapper",
			"Markdown code fence wraps the JSON payload",
			0,
			len(src),
			doc.LineIndex,
		))
	}

	if cfg.ruleEnabled("leading_or_trailing_prose") && !reMarkdownFence.MatchString(src) {
		issues = append(issues, lintFromLeadingTrailingProse(src, doc.LineIndex)...)
	}

	if cfg.ruleEnabled("single_quotes") {
		issues = append(issues, lintFromRegexMatches(
			src,
			doc.LineIndex,
			reSingleKey,
			"single_quotes",
			"Single-quoted object key used where JSON requires double quotes",
		)...)
		issues = append(issues, lintFromRegexMatches(
			src,
			doc.LineIndex,
			reSingleValue,
			"single_quotes",
			"Single-quoted string value used where JSON requires double quotes",
		)...)
	}

	if cfg.ruleEnabled("unquoted_keys") {
		issues = append(issues, lintFromRegexMatches(
			src,
			doc.LineIndex,
			reUnquotedKey,
			"unquoted_keys",
			"Object key is not quoted as a JSON string",
		)...)
	}

	if cfg.ruleEnabled("python_literals") {
		issues = append(issues, lintFromRegexMatches(
			src,
			doc.LineIndex,
			rePythonLiteral,
			"python_literals",
			"Python-style literal used where JSON requires true, false, or null",
		)...)
	}

	if cfg.ruleEnabled("comment") {
		for _, sp := range signals.commentSpans {
			issues = append(issues, newHeuristicIssue(
				"comment",
				"Comment syntax is present in JSON input",
				sp.start,
				sp.end,
				doc.LineIndex,
			))
		}
	}

	if cfg.ruleEnabled("trailing_comma") {
		for _, sp := range signals.trailingCommaSpans {
			issues = append(issues, newHeuristicIssue(
				"trailing_comma",
				"Trailing comma appears before a closing delimiter",
				sp.start,
				sp.end,
				doc.LineIndex,
			))
		}
	}

	if cfg.ruleEnabled("duplicate_comma") {
		for _, sp := range signals.duplicateCommaSpans {
			issues = append(issues, newHeuristicIssue(
				"duplicate_comma",
				"Duplicate comma appears in an object or array member list",
				sp.start,
				sp.end,
				doc.LineIndex,
			))
		}
	}

	if cfg.ruleEnabled("missing_closing_delimiter") && signals.missingCloseOpener != nil {
		op := *signals.missingCloseOpener
		delim := "brace"
		if op.ch == '[' {
			delim = "bracket"
		}
		issues = append(issues, newHeuristicIssue(
			"missing_closing_delimiter",
			"Likely missing closing "+delim+" for an open JSON container",
			op.pos,
			op.pos+1,
			doc.LineIndex,
		))
	}

	if cfg.ruleEnabled("ellipsis_or_placeholder") {
		for _, sp := range signals.ellipsisSpans {
			issues = append(issues, newHeuristicIssue(
				"ellipsis_or_placeholder",
				"Ellipsis placeholder appears inside the JSON input",
				sp.start,
				sp.end,
				doc.LineIndex,
			))
		}
		issues = append(issues, lintFromRegexMatches(
			src,
			doc.LineIndex,
			rePlaceholder,
			"ellipsis_or_placeholder",
			"Placeholder token appears inside the JSON input",
		)...)
		for _, sp := range signals.placeholderSpans {
			issues = append(issues, newHeuristicIssue(
				"ellipsis_or_placeholder",
				"Placeholder token appears inside the JSON input",
				sp.start,
				sp.end,
				doc.LineIndex,
			))
		}
	}

	return issues
}

func lintFromLeadingTrailingProse(src string, li lineIndex) []LintIssue {
	var issues []LintIssue

	trimmedStart := firstNonWhitespaceIndex(src)
	if trimmedStart == -1 {
		return nil
	}
	trimmedEnd := lastNonWhitespaceIndex(src)
	if trimmedEnd == -1 {
		return nil
	}

	firstStructural := strings.IndexAny(src, "{[")
	if firstStructural > trimmedStart {
		issues = append(issues, newHeuristicIssue(
			"leading_or_trailing_prose",
			"Non-JSON prose appears before the JSON payload",
			trimmedStart,
			firstStructural,
			li,
		))
	}

	lastStructural := strings.LastIndexAny(src, "}]")
	if lastStructural >= 0 && lastStructural < trimmedEnd {
		issues = append(issues, newHeuristicIssue(
			"leading_or_trailing_prose",
			"Non-JSON prose appears after the JSON payload",
			lastStructural+1,
			trimmedEnd+1,
			li,
		))
	}

	return issues
}

func lintFromRegexMatches(src string, li lineIndex, re *regexp.Regexp, rule, description string) []LintIssue {
	matches := re.FindAllStringIndex(src, -1)
	issues := make([]LintIssue, 0, len(matches))
	for _, match := range matches {
		issues = append(issues, newHeuristicIssue(rule, description, match[0], match[1], li))
	}
	return issues
}

func newHeuristicIssue(rule, description string, start, end int, li lineIndex) LintIssue {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	startRow, startCol := li.rowColAtByte(uint(start))
	endRow, endCol := li.rowColAtByte(uint(end))
	return LintIssue{
		Rule:        rule,
		Source:      "heuristic",
		Description: description,
		StartByte:   uint(start),
		EndByte:     uint(end),
		StartRow:    uint(startRow),
		StartCol:    uint(startCol),
		EndRow:      uint(endRow),
		EndCol:      uint(endCol),
		Row:         startRow,
	}
}

func scanHeuristicSignals(src string) heuristicSignals {
	var signals heuristicSignals
	var stack []opener

	var quote byte
	escape := false

	for i := 0; i < len(src); i++ {
		ch := src[i]

		if quote != 0 {
			if escape {
				escape = false
				continue
			}
			if ch == '\\' {
				escape = true
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}

		if ch == '"' || ch == '\'' {
			quote = ch
			continue
		}

		if ch == '/' && i+1 < len(src) {
			switch src[i+1] {
			case '/':
				end := i + 2
				for end < len(src) && src[end] != '\n' {
					end++
				}
				signals.commentSpans = append(signals.commentSpans, span{start: i, end: end})
				i = end - 1
				continue
			case '*':
				end := i + 2
				for end+1 < len(src) && (src[end] != '*' || src[end+1] != '/') {
					end++
				}
				if end+1 < len(src) {
					end += 2
				} else {
					end = len(src)
				}
				signals.commentSpans = append(signals.commentSpans, span{start: i, end: end})
				i = end - 1
				continue
			}
		}

		switch ch {
		case '{', '[':
			stack = append(stack, opener{ch: ch, pos: i})
		case '}', ']':
			if len(stack) == 0 {
				continue
			}
			last := stack[len(stack)-1]
			if (last.ch == '{' && ch == '}') || (last.ch == '[' && ch == ']') {
				stack = stack[:len(stack)-1]
			}
		case ',':
			j := skipJSONishWhitespace(src, i+1)
			if j >= len(src) {
				continue
			}
			switch src[j] {
			case ',':
				signals.duplicateCommaSpans = append(signals.duplicateCommaSpans, span{start: i, end: j + 1})
			case '}', ']':
				signals.trailingCommaSpans = append(signals.trailingCommaSpans, span{start: i, end: j + 1})
			}
		case '.':
			if i+2 < len(src) && src[i:i+3] == "..." {
				signals.ellipsisSpans = append(signals.ellipsisSpans, span{start: i, end: i + 3})
				i += 2
			}
		}
	}

	if len(stack) > 0 {
		last := stack[len(stack)-1]
		signals.missingCloseOpener = &last
	}
	return signals
}

func firstNonWhitespaceIndex(src string) int {
	for i, r := range src {
		if !unicode.IsSpace(r) {
			return i
		}
	}
	return -1
}

func lastNonWhitespaceIndex(src string) int {
	for i := len(src) - 1; i >= 0; i-- {
		if !unicode.IsSpace(rune(src[i])) {
			return i
		}
	}
	return -1
}

func skipJSONishWhitespace(src string, start int) int {
	for start < len(src) {
		if !unicode.IsSpace(rune(src[start])) {
			return start
		}
		start++
	}
	return start
}
