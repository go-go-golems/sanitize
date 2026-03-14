package jsonsanitize

import (
	"fmt"
	"sort"
	"strings"
)

// Option configures JSON sanitizer behavior.
type Option func(*config)

type config struct {
	maxIterations int
	onlyRules     map[string]bool
	disabledRules map[string]bool
}

func defaultConfig() config {
	return config{
		maxIterations: 10,
		onlyRules:     nil,
		disabledRules: map[string]bool{},
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

// WithRules restricts the sanitizer to only the named rules.
func WithRules(rules ...string) Option {
	return WithOnlyRules(rules...)
}

// WithOnlyRules restricts linting and fixing to only the named rules.
func WithOnlyRules(rules ...string) Option {
	return func(c *config) {
		c.onlyRules = make(map[string]bool, len(rules))
		for _, r := range rules {
			c.onlyRules[r] = true
		}
	}
}

// WithDisabledRules disables the named rules while leaving the others enabled.
func WithDisabledRules(rules ...string) Option {
	return func(c *config) {
		if c.disabledRules == nil {
			c.disabledRules = make(map[string]bool, len(rules))
		}
		for _, r := range rules {
			c.disabledRules[r] = true
		}
	}
}

func buildConfig(opts ...Option) (config, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if err := cfg.validate(); err != nil {
		return config{}, err
	}
	return cfg, nil
}

func (c config) validate() error {
	if err := ValidateRuleNames(keys(c.onlyRules)...); err != nil {
		return err
	}
	if err := ValidateRuleNames(keys(c.disabledRules)...); err != nil {
		return err
	}

	var overlap []string
	for name := range c.onlyRules {
		if c.disabledRules[name] {
			overlap = append(overlap, name)
		}
	}
	if len(overlap) > 0 {
		sort.Strings(overlap)
		return fmt.Errorf("rule(s) cannot be both enabled and disabled: %s", strings.Join(overlap, ", "))
	}
	return nil
}

func keys(m map[string]bool) []string {
	ret := make([]string, 0, len(m))
	for name := range m {
		ret = append(ret, name)
	}
	return ret
}
