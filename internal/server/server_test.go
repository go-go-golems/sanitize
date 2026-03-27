package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

func TestSanitizeHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/sanitize", strings.NewReader(`{"format":"yaml","input":"name:Alice\n"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	sanitizeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	var result yamlsanitize.Result
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	if result.Sanitized != "name: Alice\n" {
		t.Fatalf("expected sanitized YAML, got %q", result.Sanitized)
	}
}

func TestSanitizeHandlerRejectsOversizedBody(t *testing.T) {
	body := bytes.NewBufferString(`{"format":"yaml","input":"`)
	body.WriteString(strings.Repeat("a", maxRequestBodyBytes))
	body.WriteString(`"}`)

	req := httptest.NewRequest(http.MethodPost, "/api/sanitize", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	sanitizeHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%q", rec.Code, rec.Body.String())
	}
}

func TestSanitizeHandlerRejectsTrailingJSONValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/sanitize", strings.NewReader(`{"format":"yaml","input":"name:Alice\n"} {"extra":true}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	sanitizeHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%q", rec.Code, rec.Body.String())
	}
}

func TestParseHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/parse", nil)
	rec := httptest.NewRecorder()

	parseHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestParseHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/parse", strings.NewReader(`{"format":"yaml","input":"name: Alice\n"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	parseHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	var payload struct {
		TreeText string                   `json:"tree_text"`
		Errors   []yamlsanitize.ErrorNode `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	if payload.TreeText == "" {
		t.Fatal("expected parse tree text in response")
	}
	if len(payload.Errors) != 0 {
		t.Fatalf("expected no parse errors, got %+v", payload.Errors)
	}
}

func TestSanitizeHandlerJSONSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/sanitize", strings.NewReader("{\"format\":\"json\",\"input\":\"```json\\n{\\\"ok\\\": True,}\\n```\"}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	sanitizeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	var result jsonsanitize.Result
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	if result.Sanitized != "{\"ok\": true}\n" {
		t.Fatalf("expected sanitized JSON, got %q", result.Sanitized)
	}
}

func TestParseHandlerJSONSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/parse", strings.NewReader(`{"format":"json","input":"{\"name\":\"Alice\"}"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	parseHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	var payload struct {
		TreeText         string                   `json:"tree_text"`
		Errors           []jsonsanitize.ErrorNode `json:"errors"`
		StrictParseClean bool                     `json:"strict_parse_clean"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	if payload.TreeText == "" {
		t.Fatal("expected parse tree text in response")
	}
	if len(payload.Errors) != 0 {
		t.Fatalf("expected no parse errors, got %+v", payload.Errors)
	}
	if !payload.StrictParseClean {
		t.Fatalf("expected strict parse clean response, got %+v", payload)
	}
}

func TestParseHandlerRejectsTrailingGarbage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/parse", strings.NewReader(`{"format":"json","input":"{\"name\":\"Alice\"}"} trailing`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	parseHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%q", rec.Code, rec.Body.String())
	}
}

func TestExamplesHandlerIncludesYAMLAndJSONExamples(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/examples", nil)
	rec := httptest.NewRecorder()

	examplesHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	var payload []struct {
		Input  string `json:"input"`
		Format string `json:"format"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	var foundYAML, foundJSON bool
	for _, ex := range payload {
		if ex.Input == "" {
			t.Fatalf("expected input to be populated, got %+v", ex)
		}
		if ex.Format == "yaml" {
			foundYAML = true
		}
		if ex.Format == "json" {
			foundJSON = true
		}
	}
	if !foundYAML || !foundJSON {
		t.Fatalf("expected both yaml and json examples, got %+v", payload)
	}
}

func TestRootServesFormatAwarePlayground(t *testing.T) {
	srv, err := New(DefaultPort)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	for _, needle := range []string{
		"Sanitize Playground",
		`id="format-select"`,
		`id="example-select"`,
		`id="strict-badge"`,
	} {
		if !strings.Contains(body, needle) {
			t.Fatalf("expected %q in root HTML, got %q", needle, body)
		}
	}
}

func TestStaticAppUsesFormatAwareRequestBody(t *testing.T) {
	srv, err := New(DefaultPort)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/js/app.js", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	js := string(body)
	for _, needle := range []string{
		"JSON.stringify({ format: currentFormat, input: src })",
		"tree-sitter-json",
		"strict invalid",
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("expected %q in app.js, got %q", needle, js)
		}
	}
}
