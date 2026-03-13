package yamlsanitize

import (
	"strings"
	"testing"
)

func TestRuleCatalogIncludesKnownRules(t *testing.T) {
	names := map[string]bool{}
	for _, spec := range RuleCatalog() {
		names[spec.Name] = true
	}

	for _, name := range []string{
		"tab_indent",
		"missing_space_after_colon",
		"duplicate_key",
		"structural_parse_error",
	} {
		if !names[name] {
			t.Fatalf("expected rule catalog to contain %q", name)
		}
	}
}

func TestValidateRuleNamesRejectsUnknownRules(t *testing.T) {
	err := ValidateRuleNames("tab_indent", "not_a_rule")
	if err == nil {
		t.Fatal("expected unknown rule validation error")
	}
	if !strings.Contains(err.Error(), "not_a_rule") {
		t.Fatalf("expected error to mention unknown rule, got %v", err)
	}
}

func TestLintWithOptions_OnlyRules(t *testing.T) {
	issues, err := LintWithOptions("server:\n\thost:localhost\n", WithOnlyRules("missing_space_after_colon"))
	if err != nil {
		t.Fatalf("LintWithOptions: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected exactly one issue, got %+v", issues)
	}
	if issues[0].Rule != "missing_space_after_colon" {
		t.Fatalf("expected missing_space_after_colon, got %+v", issues)
	}
}

func TestLintWithOptions_DisabledRules(t *testing.T) {
	issues, err := LintWithOptions("server:\n\thost:localhost\n", WithDisabledRules("tab_indent"))
	if err != nil {
		t.Fatalf("LintWithOptions: %v", err)
	}
	for _, issue := range issues {
		if issue.Rule == "tab_indent" {
			t.Fatalf("did not expect tab_indent issue, got %+v", issues)
		}
	}
}

func TestLintWithOptions_UnknownRuleReturnsError(t *testing.T) {
	_, err := LintWithOptions("name:Alice\n", WithOnlyRules("not_a_rule"))
	if err == nil {
		t.Fatal("expected invalid rule selection error")
	}
}

func TestSanitizeWithOptions_DisabledRuleSkipsFix(t *testing.T) {
	result, err := SanitizeWithOptions("server:\n\thost:localhost\n", WithDisabledRules("missing_space_after_colon"))
	if err != nil {
		t.Fatalf("SanitizeWithOptions: %v", err)
	}
	if strings.Contains(result.Sanitized, "host: localhost") {
		t.Fatalf("expected missing_space_after_colon to remain disabled, got:\n%s", result.Sanitized)
	}
	if strings.Contains(result.Sanitized, "\t") {
		t.Fatalf("expected tab_indent to remain enabled, got:\n%s", result.Sanitized)
	}
}

func TestSanitizeWithOptions_UnknownRuleReturnsError(t *testing.T) {
	_, err := SanitizeWithOptions("name:Alice\n", WithOnlyRules("not_a_rule"))
	if err == nil {
		t.Fatal("expected invalid rule selection error")
	}
}

func TestSanitizeWithOptions_ConflictingRuleSelectionReturnsError(t *testing.T) {
	_, err := SanitizeWithOptions("name:Alice\n", WithOnlyRules("tab_indent"), WithDisabledRules("tab_indent"))
	if err == nil {
		t.Fatal("expected conflicting rule selection error")
	}
	if !strings.Contains(err.Error(), "tab_indent") {
		t.Fatalf("expected conflicting rule error to mention tab_indent, got %v", err)
	}
}
