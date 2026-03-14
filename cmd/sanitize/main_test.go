package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

func TestRunLintTextReturnsNonZeroOnIssues(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"lint"}, strings.NewReader("name:Alice\n"), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout.String(), "missing space after colon") {
		t.Fatalf("expected lint output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunLintJSONReturnsNonZeroOnIssues(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"lint", "--json"}, strings.NewReader("name:Alice\n"), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	var issues []yamlsanitize.LintIssue
	if err := json.Unmarshal(stdout.Bytes(), &issues); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if len(issues) != 1 || issues[0].Rule != "missing_space_after_colon" {
		t.Fatalf("expected missing_space_after_colon issue, got %+v", issues)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunSanitizeJSONReturnsZeroAfterFix(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"fix", "--json"}, strings.NewReader("name:Alice\n"), stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var result yamlsanitize.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if result.Sanitized != "name: Alice\n" {
		t.Fatalf("expected sanitized output, got %q", result.Sanitized)
	}
	if !result.ParseClean || !result.LintClean {
		t.Fatalf("expected clean result, got parse_clean=%v lint_clean=%v", result.ParseClean, result.LintClean)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunSanitizeJSONReturnsNonZeroWhenErrorsRemain(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"fix", "--json"}, strings.NewReader("a: [1,2\n"), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	var result yamlsanitize.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if result.ParseClean {
		t.Fatalf("expected parse_clean=false, got %+v", result)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunServeReturnsNonZeroOnInvalidPort(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"serve", "--port", "70000"}, strings.NewReader(""), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "invalid port") {
		t.Fatalf("expected invalid port error, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
}

func TestRunParseJSONReturnsZeroForCleanInput(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"parse", "--json"}, strings.NewReader("name: Alice\n"), stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var result struct {
		TreeText string                   `json:"tree_text"`
		Errors   []yamlsanitize.ErrorNode `json:"errors"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if result.TreeText == "" {
		t.Fatal("expected non-empty parse tree")
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected no parse errors, got %+v", result.Errors)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunParseTextReturnsNonZeroOnErrors(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"parse"}, strings.NewReader("foo: bar: baz\n"), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if stdout.Len() == 0 {
		t.Fatal("expected parse tree on stdout")
	}
	if !strings.Contains(stderr.String(), "parse error(s)") {
		t.Fatalf("expected parse summary on stderr, got %q", stderr.String())
	}
}

func TestRunLintJSONWithRuleFilter(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"lint", "--json", "--rule", "missing_space_after_colon"}, strings.NewReader("server:\n\thost:localhost\n"), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	var issues []yamlsanitize.LintIssue
	if err := json.Unmarshal(stdout.Bytes(), &issues); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if len(issues) != 1 || issues[0].Rule != "missing_space_after_colon" {
		t.Fatalf("expected only missing_space_after_colon, got %+v", issues)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunLintReturnsNonZeroOnUnknownRule(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"lint", "--rule", "not_a_rule"}, strings.NewReader("name:Alice\n"), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "invalid rule selection") {
		t.Fatalf("expected invalid rule selection error, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
}

func TestRunFixWithDisabledRuleSkipsFix(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"fix", "--disable-rule", "missing_space_after_colon"}, strings.NewReader("server:\n\thost:localhost\n"), stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0 because the remaining issue is disabled, got %d", exitCode)
	}
	if strings.Contains(stdout.String(), "host: localhost") {
		t.Fatalf("expected missing_space_after_colon to remain disabled, got %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "\t") {
		t.Fatalf("expected tab_indent to remain enabled, got %q", stdout.String())
	}
}

func TestRunRulesJSONReturnsCatalog(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"rules", "--json"}, strings.NewReader(""), stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var rules []yamlsanitize.RuleSpec
	if err := json.Unmarshal(stdout.Bytes(), &rules); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if len(rules) == 0 {
		t.Fatal("expected non-empty rule catalog")
	}
	found := false
	for _, rule := range rules {
		if rule.Name == "tab_indent" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected tab_indent in rule catalog, got %+v", rules)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunParseJSONFormatReturnsZeroForCleanInput(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"parse", "--format", "json", "--json"}, strings.NewReader(`{"name":"Alice"}`), stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var result struct {
		TreeText string                   `json:"tree_text"`
		Errors   []jsonsanitize.ErrorNode `json:"errors"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if result.TreeText == "" {
		t.Fatal("expected non-empty parse tree")
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected no parse errors, got %+v", result.Errors)
	}
}

func TestRunLintJSONFormatReturnsHeuristicIssue(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"lint", "--format", "json", "--json", "--rule", "python_literals"}, strings.NewReader(`{"ok": True}`), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	var issues []jsonsanitize.LintIssue
	if err := json.Unmarshal(stdout.Bytes(), &issues); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if len(issues) != 1 || issues[0].Rule != "python_literals" {
		t.Fatalf("expected python_literals issue, got %+v", issues)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunRulesJSONFormatReturnsCatalog(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"rules", "--format", "json", "--json"}, strings.NewReader(""), stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var rules []jsonsanitize.RuleSpec
	if err := json.Unmarshal(stdout.Bytes(), &rules); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	found := false
	for _, rule := range rules {
		if rule.Name == "python_literals" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected python_literals in rule catalog, got %+v", rules)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunFixJSONFormatReturnsNotImplemented(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run([]string{"fix", "--format", "json"}, strings.NewReader(`{"ok": True}`), stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "json fix is not implemented yet") {
		t.Fatalf("expected not implemented error, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
}
