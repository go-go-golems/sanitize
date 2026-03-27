package yamlsanitize

import (
	"fmt"
	"sort"
	"strings"
)

// RuleSpec describes one configurable lint/fix rule.
type RuleSpec struct {
	Name           string `json:"name"`
	Summary        string `json:"summary"`
	Lints          bool   `json:"lints"`
	Fixes          bool   `json:"fixes"`
	DefaultEnabled bool   `json:"default_enabled"`
}

var ruleCatalog = []RuleSpec{
	{Name: "tab_indent", Summary: "Tabs used for indentation", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "missing_space_after_colon", Summary: "Missing space after ':' in a mapping", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "list_dash_no_space", Summary: "List dash not followed by a space", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "trailing_comma", Summary: "Trailing comma in a flow collection", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "extra_colon_in_value", Summary: "Ambiguous colon in a plain scalar value", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "duplicate_key", Summary: "Duplicate key within the same mapping", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "mixed_indent", Summary: "Mixed indentation widths in a parse-damaged region", Lints: true, Fixes: true, DefaultEnabled: true},
	{Name: "structural_parse_error", Summary: "Tree-sitter structural parse error surfaced as lint", Lints: true, Fixes: false, DefaultEnabled: true},
	{Name: "missing_syntax_node", Summary: "Tree-sitter missing syntax node surfaced as lint", Lints: true, Fixes: false, DefaultEnabled: true},
}

var ruleCatalogByName = func() map[string]RuleSpec {
	m := make(map[string]RuleSpec, len(ruleCatalog))
	for _, spec := range ruleCatalog {
		m[spec.Name] = spec
	}
	return m
}()

// RuleCatalog returns a copy of the known YAML rule catalog.
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
