package yamlsanitize

import (
	"strings"
	"testing"
)

// ── ParseTree tests ──────────────────────────────────────────────────────

func TestParseTree_ValidYAML(t *testing.T) {
	src := "name: Alice\nage: 30\n"
	tree, errs, err := ParseTree(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d: %+v", len(errs), errs)
	}
	if tree == "" {
		t.Error("expected non-empty tree text")
	}
}

func TestParseTree_InvalidYAML(t *testing.T) {
	src := "foo: bar: baz\n"
	_, errs, err := ParseTree(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected at least one error for 'foo: bar: baz'")
	}
}

func TestParseTree_Empty(t *testing.T) {
	tree, errs, err := ParseTree("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors for empty input, got %d", len(errs))
	}
	if tree == "" {
		t.Error("expected non-empty tree for empty input (stream node)")
	}
}

// ── Lint tests ───────────────────────────────────────────────────────────

func TestLint_TabIndent(t *testing.T) {
	src := "server:\n\thost: localhost\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "tab_indent" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected tab_indent lint issue")
	}
}

func TestLint_MissingSpaceAfterColon(t *testing.T) {
	src := "name:Alice\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "missing_space_after_colon" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected missing_space_after_colon lint issue")
	}
}

func TestLint_MissingSpaceAfterColon_URL(t *testing.T) {
	// URLs should not be flagged
	src := "url: https://example.com\n"
	issues := Lint(src)
	for _, li := range issues {
		if li.Rule == "missing_space_after_colon" {
			t.Error("URL should not be flagged as missing_space_after_colon")
		}
	}
}

func TestLint_BadListDash(t *testing.T) {
	src := "items:\n  -apple\n  - banana\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "list_dash_no_space" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected list_dash_no_space lint issue")
	}
}

func TestLint_TrailingComma(t *testing.T) {
	src := "point: { x: 1, y: 2, }\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "trailing_comma" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected trailing_comma lint issue")
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	src := "config:\n  timeout: 30\n  timeout: 60\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "duplicate_key" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected duplicate_key lint issue")
	}
}

func TestLint_DuplicateKeyDifferentParents(t *testing.T) {
	src := "a:\n  timeout: 30\nb:\n  timeout: 60\n"
	issues := Lint(src)
	for _, li := range issues {
		if li.Rule == "duplicate_key" {
			t.Fatalf("did not expect duplicate_key issue across different mappings: %+v", li)
		}
	}
}

func TestLint_DuplicateKeyDifferentSequenceItems(t *testing.T) {
	src := "items:\n  - name: a\n    id: 1\n  - name: b\n    id: 2\n"
	issues := Lint(src)
	for _, li := range issues {
		if li.Rule == "duplicate_key" {
			t.Fatalf("did not expect duplicate_key issue across different sequence items: %+v", li)
		}
	}
}

func TestLint_ExtraColonInValue(t *testing.T) {
	src := "env: KEY: VALUE\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "extra_colon_in_value" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected extra_colon_in_value lint issue")
	}
}

func TestLint_Clean(t *testing.T) {
	src := "name: Alice\nage: 30\n"
	issues := Lint(src)
	if len(issues) != 0 {
		t.Errorf("expected no lint issues for clean YAML, got %d: %+v", len(issues), issues)
	}
}

func TestLint_ParseErrorProducesStructuralIssue(t *testing.T) {
	src := "foo: bar: baz\n"
	issues := Lint(src)
	found := false
	for _, li := range issues {
		if li.Rule == "structural_parse_error" {
			found = true
			if li.Source != "parse" {
				t.Fatalf("expected parse source, got %+v", li)
			}
			if li.StartRow != 0 {
				t.Fatalf("expected issue to start on row 0, got %+v", li)
			}
		}
	}
	if !found {
		t.Fatal("expected structural_parse_error lint issue")
	}
}

func TestLint_HeuristicIssueCarriesSpan(t *testing.T) {
	src := "name:Alice\n"
	issues := Lint(src)
	for _, li := range issues {
		if li.Rule == "missing_space_after_colon" {
			if li.Source != "heuristic" {
				t.Fatalf("expected heuristic source, got %+v", li)
			}
			if li.StartRow != 0 || li.EndRow != 0 {
				t.Fatalf("expected line-local span, got %+v", li)
			}
			if li.EndByte == 0 {
				t.Fatalf("expected non-zero end byte, got %+v", li)
			}
			return
		}
	}
	t.Fatal("expected missing_space_after_colon lint issue")
}

// ── Fix tests ────────────────────────────────────────────────────────────

func TestFix_TabIndent(t *testing.T) {
	src := "server:\n\thost: localhost\n\tport: 8080\n"
	result := Sanitize(src)
	if strings.Contains(result.Sanitized, "\t") {
		t.Error("expected tabs to be replaced")
	}
	if !strings.Contains(result.Sanitized, "  host: localhost") {
		t.Errorf("expected 2-space indent, got:\n%s", result.Sanitized)
	}
}

func TestFix_TabIndent_CustomWidth(t *testing.T) {
	src := "server:\n\thost: localhost\n"
	result := Sanitize(src, WithTabWidth(4))
	if !strings.Contains(result.Sanitized, "    host: localhost") {
		t.Errorf("expected 4-space indent, got:\n%s", result.Sanitized)
	}
}

func TestFix_MissingSpaceAfterColon(t *testing.T) {
	src := "name:Alice\nage:30\n"
	result := Sanitize(src)
	if !strings.Contains(result.Sanitized, "name: Alice") {
		t.Errorf("expected 'name: Alice', got:\n%s", result.Sanitized)
	}
}

func TestFix_BadListDash(t *testing.T) {
	src := "items:\n  -apple\n  -banana\n"
	result := Sanitize(src)
	if !strings.Contains(result.Sanitized, "- apple") {
		t.Errorf("expected '- apple', got:\n%s", result.Sanitized)
	}
}

func TestFix_TrailingComma(t *testing.T) {
	src := "point: { x: 1, y: 2, }\n"
	result := Sanitize(src)
	if strings.Contains(result.Sanitized, ",}") || strings.Contains(result.Sanitized, ", }") {
		t.Errorf("expected trailing comma removed, got:\n%s", result.Sanitized)
	}
}

func TestFix_DuplicateKey(t *testing.T) {
	src := "config:\n  timeout: 30\n  timeout: 60\n"
	result := Sanitize(src)
	if !strings.Contains(result.Sanitized, "timeout_2: 60") {
		t.Errorf("expected duplicate key to be renamed, got:\n%s", result.Sanitized)
	}
	fixFound := false
	for _, f := range result.Fixes {
		if f.Rule == "duplicate_key" {
			fixFound = true
			break
		}
	}
	if !fixFound {
		t.Error("expected duplicate_key fix to be applied")
	}
}

func TestFix_DuplicateKeyDifferentParents(t *testing.T) {
	src := "a:\n  timeout: 30\nb:\n  timeout: 60\n"
	result := Sanitize(src)
	if result.Sanitized != src {
		t.Fatalf("expected no duplicate-key rewrite across different mappings, got:\n%s", result.Sanitized)
	}
}

func TestFix_DuplicateKeyDifferentSequenceItems(t *testing.T) {
	src := "items:\n  - name: a\n    id: 1\n  - name: b\n    id: 2\n"
	result := Sanitize(src)
	if result.Sanitized != src {
		t.Fatalf("expected no duplicate-key rewrite across different sequence items, got:\n%s", result.Sanitized)
	}
}

func TestFix_ExtraColonInValue(t *testing.T) {
	src := "foobar: foba: sldkjf\n"
	result := Sanitize(src)
	// The value should be quoted
	if !strings.Contains(result.Sanitized, `"foba: sldkjf"`) {
		t.Errorf("expected quoted value, got:\n%s", result.Sanitized)
	}
}

// ── Sanitize integration tests ───────────────────────────────────────────

func TestSanitize_AlreadyClean(t *testing.T) {
	src := "name: Alice\nage: 30\n"
	result := Sanitize(src)
	if !result.ParseClean {
		t.Error("expected ParseClean for valid YAML")
	}
	if !result.LintClean {
		t.Error("expected LintClean for valid YAML")
	}
	if result.Sanitized != src {
		t.Errorf("expected no changes, got:\n%s", result.Sanitized)
	}
	if len(result.Fixes) != 0 {
		t.Errorf("expected no fixes, got %d", len(result.Fixes))
	}
}

func TestSanitize_MultipleErrors(t *testing.T) {
	src := "service:\n\tname:my-app\n\tport: 8080\n\ttags:\n\t  -web\n\t  -backend\n\tenv: KEY: VALUE\n"
	result := Sanitize(src)
	if len(result.Fixes) == 0 {
		t.Error("expected at least one fix for the multiple-error example")
	}
	// Should not contain any tabs after sanitization
	if strings.Contains(result.Sanitized, "\t") {
		t.Error("expected all tabs to be replaced")
	}
}

func TestSanitize_WithMaxIterations(t *testing.T) {
	src := "name:Alice\n"
	result := Sanitize(src, WithMaxIterations(1))
	// Even with 1 iteration, should fix the missing space
	if !strings.Contains(result.Sanitized, "name: Alice") {
		t.Errorf("expected fix with 1 iteration, got:\n%s", result.Sanitized)
	}
}

func TestSanitize_WithRules(t *testing.T) {
	src := "server:\n\thost:localhost\n"
	// Only fix tabs, not missing space
	result := Sanitize(src, WithRules("tab_indent"))
	if strings.Contains(result.Sanitized, "\t") {
		t.Error("expected tabs to be fixed")
	}
	// Missing space should NOT be fixed since rule is disabled
	if strings.Contains(result.Sanitized, "host: localhost") {
		t.Error("expected missing_space_after_colon to NOT be fixed with WithRules(tab_indent)")
	}
}

func TestSanitize_Empty(t *testing.T) {
	result := Sanitize("")
	if !result.ParseClean {
		t.Error("expected ParseClean for empty input")
	}
	if !result.LintClean {
		t.Error("expected LintClean for empty input")
	}
}

// ── Examples integration ─────────────────────────────────────────────────

func TestExamples_AllParse(t *testing.T) {
	for _, ex := range Examples {
		t.Run(ex.Name, func(t *testing.T) {
			result := Sanitize(ex.YAML)
			if result.Sanitized == "" && ex.YAML != "" {
				t.Error("expected non-empty sanitized output")
			}
		})
	}
}

// ── Helper unit tests ────────────────────────────────────────────────────

func TestReplaceLeadingTabs(t *testing.T) {
	tests := []struct {
		in       string
		tabWidth int
		want     string
	}{
		{"\thello", 2, "  hello"},
		{"\t\thello", 2, "    hello"},
		{"\thello", 4, "    hello"},
		{"hello", 2, "hello"},
		{"  hello", 2, "  hello"},
	}
	for _, tt := range tests {
		got := replaceLeadingTabs(tt.in, tt.tabWidth)
		if got != tt.want {
			t.Errorf("replaceLeadingTabs(%q, %d) = %q, want %q", tt.in, tt.tabWidth, got, tt.want)
		}
	}
}

func TestQuoteValue(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`hello: world`, `"hello: world"`},
		{`already "quoted"`, `"already \"quoted\""`},
		{`"pre-quoted"`, `"pre-quoted"`},
		{`'single-quoted'`, `'single-quoted'`},
	}
	for _, tt := range tests {
		got := quoteValue(tt.in)
		if got != tt.want {
			t.Errorf("quoteValue(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{4, 2, 2},
		{6, 4, 2},
		{0, 5, 5},
		{7, 7, 7},
		{12, 8, 4},
	}
	for _, tt := range tests {
		got := gcd(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("gcd(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestLineIndexAtByte(t *testing.T) {
	src := "alpha\nbeta\ngamma"
	index := newLineIndex(src)

	tests := []struct {
		offset uint
		row    int
		col    int
	}{
		{0, 0, 0},
		{2, 0, 2},
		{6, 1, 0},
		{9, 1, 3},
		{11, 2, 0},
		{15, 2, 4},
	}

	for _, tt := range tests {
		row, col := index.rowColAtByte(tt.offset)
		if row != tt.row || col != tt.col {
			t.Fatalf("rowColAtByte(%d) = (%d, %d), want (%d, %d)", tt.offset, row, col, tt.row, tt.col)
		}
	}
}
