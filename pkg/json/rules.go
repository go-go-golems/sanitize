package jsonsanitize

import (
	"fmt"
	"sort"
	"strings"
)

// RuleSpec describes one configurable JSON lint/fix rule.
type RuleSpec struct {
	Name           string `json:"name"`
	Summary        string `json:"summary"`
	Lints          bool   `json:"lints"`
	Fixes          bool   `json:"fixes"`
	ParseAware     bool   `json:"parse_aware"`
	DefaultEnabled bool   `json:"default_enabled"`
}

var ruleCatalog = []RuleSpec{
	{Name: "markdown_fence_wrapper", Summary: "Markdown code fences wrapped around the JSON payload", Lints: true, Fixes: true, ParseAware: false, DefaultEnabled: true},
	{Name: "leading_or_trailing_prose", Summary: "Non-JSON prose before or after the payload", Lints: true, Fixes: false, ParseAware: false, DefaultEnabled: true},
	{Name: "single_quotes", Summary: "Single-quoted strings used where JSON requires double quotes", Lints: true, Fixes: true, ParseAware: false, DefaultEnabled: true},
	{Name: "unquoted_keys", Summary: "Object keys are not quoted JSON strings", Lints: true, Fixes: false, ParseAware: true, DefaultEnabled: true},
	{Name: "python_literals", Summary: "Python-style literals such as True, False, and None", Lints: true, Fixes: true, ParseAware: false, DefaultEnabled: true},
	{Name: "trailing_comma", Summary: "Trailing comma in an object or array", Lints: true, Fixes: true, ParseAware: true, DefaultEnabled: true},
	{Name: "duplicate_comma", Summary: "Duplicate comma in an array or object member list", Lints: true, Fixes: true, ParseAware: false, DefaultEnabled: true},
	{Name: "comment", Summary: "Comment syntax present in JSON input", Lints: true, Fixes: true, ParseAware: false, DefaultEnabled: true},
	{Name: "missing_closing_delimiter", Summary: "Likely missing closing brace or bracket", Lints: true, Fixes: false, ParseAware: true, DefaultEnabled: true},
	{Name: "duplicate_key", Summary: "Duplicate key within the same object", Lints: true, Fixes: false, ParseAware: false, DefaultEnabled: true},
	{Name: "structural_parse_error", Summary: "Tree-sitter structural parse error surfaced as lint", Lints: true, Fixes: false, ParseAware: true, DefaultEnabled: true},
	{Name: "missing_syntax_node", Summary: "Tree-sitter missing syntax node surfaced as lint", Lints: true, Fixes: false, ParseAware: true, DefaultEnabled: true},
}

var ruleCatalogByName = func() map[string]RuleSpec {
	m := make(map[string]RuleSpec, len(ruleCatalog))
	for _, spec := range ruleCatalog {
		m[spec.Name] = spec
	}
	return m
}()

// RuleCatalog returns a copy of the known JSON rule catalog.
func RuleCatalog() []RuleSpec {
	ret := make([]RuleSpec, len(ruleCatalog))
	copy(ret, ruleCatalog)
	return ret
}

// KnownRule reports whether the provided rule name is defined in the catalog.
func KnownRule(name string) bool {
	_, ok := ruleCatalogByName[name]
	return ok
}

// LookupRule returns the named rule and whether it exists.
func LookupRule(name string) (RuleSpec, bool) {
	spec, ok := ruleCatalogByName[name]
	return spec, ok
}

// ValidateRuleNames rejects unknown rule names.
func ValidateRuleNames(names ...string) error {
	unknown := map[string]bool{}
	for _, name := range names {
		if !KnownRule(name) {
			unknown[name] = true
		}
	}
	if len(unknown) == 0 {
		return nil
	}

	ordered := make([]string, 0, len(unknown))
	for name := range unknown {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)
	return fmt.Errorf("unknown rule name(s): %s", strings.Join(ordered, ", "))
}
