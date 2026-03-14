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
