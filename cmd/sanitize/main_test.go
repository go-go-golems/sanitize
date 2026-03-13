package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

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
