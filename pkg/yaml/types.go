// Package yamlsanitize provides YAML error detection, linting, and heuristic
// fixing using tree-sitter.
//
// Strategy:
//   - Parse the input with tree-sitter (tree-sitter-grammars/tree-sitter-yaml v0.7.2).
//   - Collect all ERROR and MISSING nodes (structural parse failures).
//   - Run lint rules that detect common YAML mistakes even when tree-sitter
//     does not raise an ERROR (e.g. duplicate keys, missing space after colon,
//     bad list dash, trailing comma, tab indent).
//   - Apply fixes iteratively until the tree is clean or no more progress.
package yamlsanitize

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
	Original           string      `json:"original"`
	Sanitized          string      `json:"sanitized"`
	TreeText           string      `json:"tree_text"`
	OriginalTreeText   string      `json:"original_tree_text"`
	Errors             []ErrorNode `json:"errors"`
	OriginalErrors     []ErrorNode `json:"original_errors"`
	LintIssues         []LintIssue `json:"lint_issues"`
	OriginalLintIssues []LintIssue `json:"original_lint_issues"`
	Fixes              []Fix       `json:"fixes"`
	ParseClean         bool        `json:"parse_clean"`
	LintClean          bool        `json:"lint_clean"`
}

// Example is a named YAML snippet with a description of the error it contains.
type Example struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	YAML        string `json:"yaml"`
}
