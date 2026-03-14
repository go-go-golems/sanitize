package jsonsanitize

import (
	"strings"
	"testing"
)

func TestParseTreeCleanJSON(t *testing.T) {
	treeText, errors, err := ParseTree(`{"name":"alice","count":2}`)
	if err != nil {
		t.Fatalf("ParseTree: %v", err)
	}
	if treeText == "" {
		t.Fatal("expected non-empty tree text")
	}
	if len(errors) != 0 {
		t.Fatalf("expected no parse errors, got %+v", errors)
	}
}

func TestParseTreeMalformedJSONReturnsErrors(t *testing.T) {
	_, errors, err := ParseTree(`{"name":"alice",}`)
	if err != nil {
		t.Fatalf("ParseTree: %v", err)
	}
	if len(errors) == 0 {
		t.Fatal("expected parse errors for trailing comma JSON")
	}
}

func TestStrictParseRejectsTrailingComma(t *testing.T) {
	err := StrictParse(`{"name":"alice",}`)
	if err == nil {
		t.Fatal("expected strict parse error")
	}
	if !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("expected invalid character error, got %v", err)
	}
}

func TestStrictParseRejectsMultipleTopLevelValues(t *testing.T) {
	err := StrictParse(`{"a":1} {"b":2}`)
	if err == nil {
		t.Fatal("expected strict parse error")
	}
	if !strings.Contains(err.Error(), "multiple top-level") {
		t.Fatalf("expected multiple top-level error, got %v", err)
	}
}

func TestAnalyzeDocumentCollectsDuplicateKeys(t *testing.T) {
	doc, err := analyzeDocument("{\"a\":1,\"b\":2,\"a\":3}")
	if err != nil {
		t.Fatalf("analyzeDocument: %v", err)
	}
	if doc.StrictParseError != nil {
		t.Fatalf("expected duplicate keys to remain strict-parse valid, got %v", doc.StrictParseError)
	}
	if len(doc.DuplicateKeys) != 1 {
		t.Fatalf("expected one duplicate key, got %+v", doc.DuplicateKeys)
	}
	if doc.DuplicateKeys[0].Key != "a" {
		t.Fatalf("expected duplicate key a, got %+v", doc.DuplicateKeys[0])
	}
}

func TestAnalyzeDocumentCapturesStrictParseErrorSeparately(t *testing.T) {
	doc, err := analyzeDocument("{'a': 1}")
	if err != nil {
		t.Fatalf("analyzeDocument: %v", err)
	}
	if doc.StrictParseError == nil {
		t.Fatal("expected strict parse error for single-quoted JSON")
	}
	if len(doc.ParseErrors) == 0 {
		t.Fatal("expected tree-sitter parse errors for single-quoted JSON")
	}
}

func TestLintWithOptionsReportsDuplicateKey(t *testing.T) {
	issues, err := LintWithOptions(`{"a":1,"a":2}`)
	if err != nil {
		t.Fatalf("LintWithOptions: %v", err)
	}

	var found bool
	for _, issue := range issues {
		if issue.Rule == "duplicate_key" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected duplicate_key issue, got %+v", issues)
	}
}

func TestLintWithOptionsReportsStrictParseError(t *testing.T) {
	issues, err := LintWithOptions(`{"a":1} {"b":2}`)
	if err != nil {
		t.Fatalf("LintWithOptions: %v", err)
	}

	var found bool
	for _, issue := range issues {
		if issue.Rule == "strict_parse_error" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected strict_parse_error issue, got %+v", issues)
	}
}

func TestLintWithOptionsRuleFilter(t *testing.T) {
	issues, err := LintWithOptions(`{"a":1,"a":2}`, WithOnlyRules("duplicate_key"))
	if err != nil {
		t.Fatalf("LintWithOptions: %v", err)
	}
	if len(issues) != 1 || issues[0].Rule != "duplicate_key" {
		t.Fatalf("expected only duplicate_key, got %+v", issues)
	}
}
