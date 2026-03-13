package yamlsanitize

// Option configures the sanitizer behavior.
type Option func(*config)

type config struct {
	maxIterations int
	tabWidth      int
	enabledRules  map[string]bool // nil means all enabled
}

func defaultConfig() config {
	return config{
		maxIterations: 10,
		tabWidth:      2,
		enabledRules:  nil,
	}
}

// WithMaxIterations sets the maximum number of fix iterations (default 10).
func WithMaxIterations(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.maxIterations = n
		}
	}
}

// WithTabWidth sets the number of spaces per tab for the tab_indent fixer (default 2).
func WithTabWidth(w int) Option {
	return func(c *config) {
		if w > 0 {
			c.tabWidth = w
		}
	}
}

// WithRules restricts the sanitizer to only the named rules.
// If not called, all rules are enabled.
func WithRules(rules ...string) Option {
	return func(c *config) {
		c.enabledRules = make(map[string]bool, len(rules))
		for _, r := range rules {
			c.enabledRules[r] = true
		}
	}
}

func (c *config) ruleEnabled(rule string) bool {
	if c.enabledRules == nil {
		return true
	}
	return c.enabledRules[rule]
}
