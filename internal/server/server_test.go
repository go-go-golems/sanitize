package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

func TestSanitizeHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/sanitize", strings.NewReader(`{"yaml":"name:Alice\n"}`))
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
	body := bytes.NewBufferString(`{"yaml":"`)
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

func TestParseHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/parse", nil)
	rec := httptest.NewRecorder()

	parseHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestParseHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/parse", strings.NewReader(`{"yaml":"name: Alice\n"}`))
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
