package jsonsanitize

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
		"markdown_fence_wrapper",
		"python_literals",
		"duplicate_key",
		"structural_parse_error",
	} {
		if !names[name] {
			t.Fatalf("expected rule catalog to contain %q", name)
		}
	}
}

func TestValidateRuleNamesRejectsUnknownRules(t *testing.T) {
	err := ValidateRuleNames("python_literals", "not_a_rule")
	if err == nil {
		t.Fatal("expected unknown rule validation error")
	}
	if !strings.Contains(err.Error(), "not_a_rule") {
		t.Fatalf("expected error to mention unknown rule, got %v", err)
	}
}

func TestBuildConfigRejectsConflictingRuleSelection(t *testing.T) {
	_, err := buildConfig(WithOnlyRules("python_literals"), WithDisabledRules("python_literals"))
	if err == nil {
		t.Fatal("expected conflicting rule selection error")
	}
	if !strings.Contains(err.Error(), "python_literals") {
		t.Fatalf("expected conflicting rule error to mention python_literals, got %v", err)
	}
}
