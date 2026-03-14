// Package jsonsanitize provides JSON parse analysis, linting, and conservative
// repair helpers built on tree-sitter and strict encoding/json validation.
package jsonsanitize

// ErrorNode describes a tree-sitter ERROR or MISSING node found during parsing.
type ErrorNode struct {
	Type      string `json:"type"` // "ERROR" or "MISSING"
	StartByte uint   `json:"start_byte"`
	EndByte   uint   `json:"end_byte"`
	StartRow  uint   `json:"start_row"`
	StartCol  uint   `json:"start_col"`
	EndRow    uint   `json:"end_row"`
	EndCol    uint   `json:"end_col"`
	Text      string `json:"text"`
}

// LintIssue describes a JSON issue found by the linter.
type LintIssue struct {
	Rule        string `json:"rule"`
	Source      string `json:"source"` // parse, strict-parser, or heuristic
	Description string `json:"description"`
	StartByte   uint   `json:"start_byte"`
	EndByte     uint   `json:"end_byte"`
	StartRow    uint   `json:"start_row"`
	StartCol    uint   `json:"start_col"`
	EndRow      uint   `json:"end_row"`
	EndCol      uint   `json:"end_col"`
	Row         int    `json:"row"` // 0-indexed convenience alias for StartRow
}

// Fix describes one applied JSON fix.
type Fix struct {
	Rule        string `json:"rule"`
	Description string `json:"description"`
	Before      string `json:"before"`
	After       string `json:"after"`
}

// Result is the full output of a JSON sanitize pass.
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
	StrictParseClean   bool        `json:"strict_parse_clean"`
}

// Example is a named JSON snippet with a description of the error it contains.
type Example struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	JSON        string `json:"json"`
	Category    string `json:"category,omitempty"`
	Source      string `json:"source,omitempty"`
	Filename    string `json:"filename,omitempty"`
}
