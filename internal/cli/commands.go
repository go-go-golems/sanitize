package sanitizecli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	glazecmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/sanitize/internal/server"
	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

type fixSettings struct {
	Input         string   `glazed:"input"`
	Format        string   `glazed:"format"`
	JSON          bool     `glazed:"json"`
	TabWidth      int      `glazed:"tab-width"`
	MaxIterations int      `glazed:"max-iterations"`
	Rules         []string `glazed:"rule"`
	DisableRules  []string `glazed:"disable-rule"`
}

type lintSettings struct {
	Input        string   `glazed:"input"`
	Format       string   `glazed:"format"`
	JSON         bool     `glazed:"json"`
	Rules        []string `glazed:"rule"`
	DisableRules []string `glazed:"disable-rule"`
}

type parseSettings struct {
	Input  string `glazed:"input"`
	Format string `glazed:"format"`
	JSON   bool   `glazed:"json"`
}

type rulesSettings struct {
	Format string `glazed:"format"`
	JSON   bool   `glazed:"json"`
}

type serveSettings struct {
	Port int `glazed:"port"`
}

type fixCommand struct {
	*glazecmds.CommandDescription
	streams Streams
}

type lintCommand struct {
	*glazecmds.CommandDescription
	streams Streams
}

type parseCommand struct {
	*glazecmds.CommandDescription
	streams Streams
}

type rulesCommand struct {
	*glazecmds.CommandDescription
	streams Streams
}

type serveCommand struct {
	*glazecmds.CommandDescription
}

type parseOutput struct {
	TreeText string                   `json:"tree_text"`
	Errors   []yamlsanitize.ErrorNode `json:"errors"`
}

func newFixCommand(streams Streams) (glazecmds.Command, error) {
	sections, err := defaultSections()
	if err != nil {
		return nil, err
	}

	return &fixCommand{
		CommandDescription: glazecmds.NewCommandDescription(
			"fix",
			glazecmds.WithShort("Sanitize YAML input"),
			glazecmds.WithArguments(
				fields.New("input", fields.TypeString, fields.WithHelp("Optional input file path; reads stdin when omitted")),
			),
			glazecmds.WithFlags(
				fields.New("format", fields.TypeString, fields.WithDefault("yaml"), fields.WithHelp("Input format: yaml or json")),
				fields.New("json", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Output the full sanitize result as JSON")),
				fields.New("tab-width", fields.TypeInteger, fields.WithDefault(2), fields.WithHelp("Spaces per tab for the tab_indent fixer")),
				fields.New("max-iterations", fields.TypeInteger, fields.WithDefault(10), fields.WithHelp("Maximum number of fix iterations")),
				fields.New("rule", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Enable only the named rules (repeat with --rule)")),
				fields.New("disable-rule", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Disable the named rules (repeat with --disable-rule)")),
			),
			glazecmds.WithSections(sections...),
		),
		streams: streams,
	}, nil
}

func (c *fixCommand) Run(_ context.Context, vals *values.Values) error {
	settings := &fixSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return newExitError(1, fmt.Errorf("decode fix settings: %w", err))
	}

	format, err := normalizeFormat(settings.Format)
	if err != nil {
		return newExitError(1, err)
	}

	input, err := readInput(settings.Input, c.streams.Stdin)
	if err != nil {
		return newExitError(1, fmt.Errorf("error reading input: %w", err))
	}

	switch format {
	case "yaml":
		opts := append(
			[]yamlsanitize.Option{
				yamlsanitize.WithTabWidth(settings.TabWidth),
				yamlsanitize.WithMaxIterations(settings.MaxIterations),
			},
			buildYAMLRuleOptions(settings.Rules, settings.DisableRules)...,
		)

		result, err := yamlsanitize.SanitizeWithOptions(string(input), opts...)
		if err != nil {
			return newExitError(1, fmt.Errorf("invalid rule selection: %w", err))
		}

		if settings.JSON {
			enc := json.NewEncoder(c.streams.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(result); err != nil {
				return newExitError(1, fmt.Errorf("error encoding sanitize result: %w", err))
			}
		} else {
			if _, err := io.WriteString(c.streams.Stdout, result.Sanitized); err != nil {
				return newExitError(1, fmt.Errorf("error writing sanitized output: %w", err))
			}
			if err := writeYAMLFixSummary(c.streams.Stderr, result.Fixes); err != nil {
				return newExitError(1, err)
			}
		}

		if !result.ParseClean || !result.LintClean {
			return newExitError(1, nil)
		}
	case "json":
		opts := append(
			[]jsonsanitize.Option{
				jsonsanitize.WithMaxIterations(settings.MaxIterations),
			},
			buildJSONRuleOptions(settings.Rules, settings.DisableRules)...,
		)

		result, err := jsonsanitize.SanitizeWithOptions(string(input), opts...)
		if err != nil {
			return newExitError(1, fmt.Errorf("invalid rule selection: %w", err))
		}

		if settings.JSON {
			enc := json.NewEncoder(c.streams.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(result); err != nil {
				return newExitError(1, fmt.Errorf("error encoding sanitize result: %w", err))
			}
		} else {
			if _, err := io.WriteString(c.streams.Stdout, result.Sanitized); err != nil {
				return newExitError(1, fmt.Errorf("error writing sanitized output: %w", err))
			}
			if err := writeJSONFixSummary(c.streams.Stderr, result.Fixes); err != nil {
				return newExitError(1, err)
			}
		}

		if !result.ParseClean || !result.LintClean || !result.StrictParseClean {
			return newExitError(1, nil)
		}
	}
	return nil
}

func newLintCommand(streams Streams) (glazecmds.Command, error) {
	sections, err := defaultSections()
	if err != nil {
		return nil, err
	}

	return &lintCommand{
		CommandDescription: glazecmds.NewCommandDescription(
			"lint",
			glazecmds.WithShort("Lint YAML input without applying fixes"),
			glazecmds.WithArguments(
				fields.New("input", fields.TypeString, fields.WithHelp("Optional input file path; reads stdin when omitted")),
			),
			glazecmds.WithFlags(
				fields.New("format", fields.TypeString, fields.WithDefault("yaml"), fields.WithHelp("Input format: yaml or json")),
				fields.New("json", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Output lint issues as JSON")),
				fields.New("rule", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Enable only the named rules (repeat with --rule)")),
				fields.New("disable-rule", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Disable the named rules (repeat with --disable-rule)")),
			),
			glazecmds.WithSections(sections...),
		),
		streams: streams,
	}, nil
}

func (c *lintCommand) Run(_ context.Context, vals *values.Values) error {
	settings := &lintSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return newExitError(1, fmt.Errorf("decode lint settings: %w", err))
	}

	format, err := normalizeFormat(settings.Format)
	if err != nil {
		return newExitError(1, err)
	}

	input, err := readInput(settings.Input, c.streams.Stdin)
	if err != nil {
		return newExitError(1, fmt.Errorf("error reading input: %w", err))
	}

	switch format {
	case "yaml":
		issues, err := yamlsanitize.LintWithOptions(string(input), buildYAMLRuleOptions(settings.Rules, settings.DisableRules)...)
		if err != nil {
			return newExitError(1, fmt.Errorf("invalid rule selection: %w", err))
		}
		if settings.JSON {
			if err := json.NewEncoder(c.streams.Stdout).Encode(issues); err != nil {
				return newExitError(1, fmt.Errorf("error encoding lint result: %w", err))
			}
		} else {
			for _, issue := range issues {
				if _, err := fmt.Fprintln(c.streams.Stdout, issue.Description); err != nil {
					return newExitError(1, fmt.Errorf("error writing lint output: %w", err))
				}
			}
		}
		if len(issues) > 0 {
			return newExitError(1, nil)
		}
	case "json":
		issues, err := jsonsanitize.LintWithOptions(string(input), buildJSONRuleOptions(settings.Rules, settings.DisableRules)...)
		if err != nil {
			return newExitError(1, fmt.Errorf("invalid rule selection: %w", err))
		}
		if settings.JSON {
			if err := json.NewEncoder(c.streams.Stdout).Encode(issues); err != nil {
				return newExitError(1, fmt.Errorf("error encoding lint result: %w", err))
			}
		} else {
			for _, issue := range issues {
				if _, err := fmt.Fprintln(c.streams.Stdout, issue.Description); err != nil {
					return newExitError(1, fmt.Errorf("error writing lint output: %w", err))
				}
			}
		}
		if len(issues) > 0 {
			return newExitError(1, nil)
		}
	}
	return nil
}

func newRulesCommand(streams Streams) (glazecmds.Command, error) {
	sections, err := defaultSections()
	if err != nil {
		return nil, err
	}

	return &rulesCommand{
		CommandDescription: glazecmds.NewCommandDescription(
			"rules",
			glazecmds.WithShort("List the available YAML lint and fix rules"),
			glazecmds.WithFlags(
				fields.New("format", fields.TypeString, fields.WithDefault("yaml"), fields.WithHelp("Input format: yaml or json")),
				fields.New("json", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Output rule metadata as JSON")),
			),
			glazecmds.WithSections(sections...),
		),
		streams: streams,
	}, nil
}

func (c *rulesCommand) Run(_ context.Context, vals *values.Values) error {
	settings := &rulesSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return newExitError(1, fmt.Errorf("decode rules settings: %w", err))
	}

	format, err := normalizeFormat(settings.Format)
	if err != nil {
		return newExitError(1, err)
	}

	if settings.JSON {
		enc := json.NewEncoder(c.streams.Stdout)
		enc.SetIndent("", "  ")
		switch format {
		case "yaml":
			if err := enc.Encode(yamlsanitize.RuleCatalog()); err != nil {
				return newExitError(1, fmt.Errorf("error encoding rule list: %w", err))
			}
		case "json":
			if err := enc.Encode(jsonsanitize.RuleCatalog()); err != nil {
				return newExitError(1, fmt.Errorf("error encoding rule list: %w", err))
			}
		}
		return nil
	}

	tw := tabwriter.NewWriter(c.streams.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "NAME\tLINTS\tFIXES\tDEFAULT\tSUMMARY"); err != nil {
		return newExitError(1, fmt.Errorf("error writing rule header: %w", err))
	}
	switch format {
	case "yaml":
		for _, rule := range yamlsanitize.RuleCatalog() {
			if _, err := fmt.Fprintf(tw, "%s\t%t\t%t\t%t\t%s\n", rule.Name, rule.Lints, rule.Fixes, rule.DefaultEnabled, rule.Summary); err != nil {
				return newExitError(1, fmt.Errorf("error writing rule list: %w", err))
			}
		}
	case "json":
		for _, rule := range jsonsanitize.RuleCatalog() {
			if _, err := fmt.Fprintf(tw, "%s\t%t\t%t\t%t\t%s\n", rule.Name, rule.Lints, rule.Fixes, rule.DefaultEnabled, rule.Summary); err != nil {
				return newExitError(1, fmt.Errorf("error writing rule list: %w", err))
			}
		}
	}
	if err := tw.Flush(); err != nil {
		return newExitError(1, fmt.Errorf("error flushing rule list: %w", err))
	}
	return nil
}

func newParseCommand(streams Streams) (glazecmds.Command, error) {
	sections, err := defaultSections()
	if err != nil {
		return nil, err
	}

	return &parseCommand{
		CommandDescription: glazecmds.NewCommandDescription(
			"parse",
			glazecmds.WithShort("Print the tree-sitter parse tree and structural errors"),
			glazecmds.WithArguments(
				fields.New("input", fields.TypeString, fields.WithHelp("Optional input file path; reads stdin when omitted")),
			),
			glazecmds.WithFlags(
				fields.New("format", fields.TypeString, fields.WithDefault("yaml"), fields.WithHelp("Input format: yaml or json")),
				fields.New("json", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Output parse tree and errors as JSON")),
			),
			glazecmds.WithSections(sections...),
		),
		streams: streams,
	}, nil
}

func (c *parseCommand) Run(_ context.Context, vals *values.Values) error {
	settings := &parseSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return newExitError(1, fmt.Errorf("decode parse settings: %w", err))
	}

	format, err := normalizeFormat(settings.Format)
	if err != nil {
		return newExitError(1, err)
	}

	input, err := readInput(settings.Input, c.streams.Stdin)
	if err != nil {
		return newExitError(1, fmt.Errorf("error reading input: %w", err))
	}

	switch format {
	case "yaml":
		treeText, errors, err := yamlsanitize.ParseTree(string(input))
		if err != nil {
			return newExitError(1, fmt.Errorf("error parsing input: %w", err))
		}
		result := parseOutput{
			TreeText: treeText,
			Errors:   errors,
		}
		if settings.JSON {
			enc := json.NewEncoder(c.streams.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(result); err != nil {
				return newExitError(1, fmt.Errorf("error encoding parse result: %w", err))
			}
		} else {
			if _, err := fmt.Fprintln(c.streams.Stdout, result.TreeText); err != nil {
				return newExitError(1, fmt.Errorf("error writing parse tree: %w", err))
			}
			if err := writeYAMLParseSummary(c.streams.Stderr, result.Errors); err != nil {
				return newExitError(1, err)
			}
		}
		if len(result.Errors) > 0 {
			return newExitError(1, nil)
		}
	case "json":
		treeText, errors, err := jsonsanitize.ParseTree(string(input))
		if err != nil {
			return newExitError(1, fmt.Errorf("error parsing input: %w", err))
		}
		result := struct {
			TreeText string                   `json:"tree_text"`
			Errors   []jsonsanitize.ErrorNode `json:"errors"`
		}{
			TreeText: treeText,
			Errors:   errors,
		}
		if settings.JSON {
			enc := json.NewEncoder(c.streams.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(result); err != nil {
				return newExitError(1, fmt.Errorf("error encoding parse result: %w", err))
			}
		} else {
			if _, err := fmt.Fprintln(c.streams.Stdout, result.TreeText); err != nil {
				return newExitError(1, fmt.Errorf("error writing parse tree: %w", err))
			}
			if err := writeJSONParseSummary(c.streams.Stderr, result.Errors); err != nil {
				return newExitError(1, err)
			}
		}
		if len(result.Errors) > 0 {
			return newExitError(1, nil)
		}
	}
	return nil
}

func newServeCommand(_ Streams) (glazecmds.Command, error) {
	sections, err := defaultSections()
	if err != nil {
		return nil, err
	}

	return &serveCommand{
		CommandDescription: glazecmds.NewCommandDescription(
			"serve",
			glazecmds.WithShort("Run the sanitize web server"),
			glazecmds.WithFlags(
				fields.New("port", fields.TypeInteger, fields.WithDefault(server.DefaultPort), fields.WithHelp("HTTP port for the web UI and API")),
			),
			glazecmds.WithSections(sections...),
		),
	}, nil
}

func (c *serveCommand) Run(_ context.Context, vals *values.Values) error {
	settings := &serveSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return newExitError(1, fmt.Errorf("decode serve settings: %w", err))
	}

	if err := server.Run(settings.Port); err != nil {
		return newExitError(1, err)
	}
	return nil
}

func buildYAMLRuleOptions(rules, disabledRules []string) []yamlsanitize.Option {
	opts := make([]yamlsanitize.Option, 0, 2)
	if len(rules) > 0 {
		opts = append(opts, yamlsanitize.WithOnlyRules(rules...))
	}
	if len(disabledRules) > 0 {
		opts = append(opts, yamlsanitize.WithDisabledRules(disabledRules...))
	}
	return opts
}

func buildJSONRuleOptions(rules, disabledRules []string) []jsonsanitize.Option {
	opts := make([]jsonsanitize.Option, 0, 2)
	if len(rules) > 0 {
		opts = append(opts, jsonsanitize.WithOnlyRules(rules...))
	}
	if len(disabledRules) > 0 {
		opts = append(opts, jsonsanitize.WithDisabledRules(disabledRules...))
	}
	return opts
}

func readInput(path string, stdin io.Reader) ([]byte, error) {
	if path != "" {
		return os.ReadFile(path)
	}
	return io.ReadAll(stdin)
}

func writeYAMLFixSummary(w io.Writer, fixes []yamlsanitize.Fix) error {
	if len(fixes) == 0 {
		return nil
	}

	if _, err := fmt.Fprintf(w, "%d fix(es) applied\n", len(fixes)); err != nil {
		return fmt.Errorf("error writing fix summary: %w", err)
	}
	for _, fix := range fixes {
		if _, err := fmt.Fprintf(w, "  %s: %s\n", fix.Rule, fix.Description); err != nil {
			return fmt.Errorf("error writing fix summary: %w", err)
		}
	}
	return nil
}

func writeJSONFixSummary(w io.Writer, fixes []jsonsanitize.Fix) error {
	if len(fixes) == 0 {
		return nil
	}

	if _, err := fmt.Fprintf(w, "%d fix(es) applied\n", len(fixes)); err != nil {
		return fmt.Errorf("error writing fix summary: %w", err)
	}
	for _, fix := range fixes {
		if _, err := fmt.Fprintf(w, "  %s: %s\n", fix.Rule, fix.Description); err != nil {
			return fmt.Errorf("error writing fix summary: %w", err)
		}
	}
	return nil
}

func writeYAMLParseSummary(w io.Writer, errors []yamlsanitize.ErrorNode) error {
	if len(errors) == 0 {
		if _, err := io.WriteString(w, "0 parse error(s)\n"); err != nil {
			return fmt.Errorf("error writing parse summary: %w", err)
		}
		return nil
	}

	if _, err := fmt.Fprintf(w, "%d parse error(s)\n", len(errors)); err != nil {
		return fmt.Errorf("error writing parse summary: %w", err)
	}
	for _, parseErr := range errors {
		if _, err := fmt.Fprintf(
			w,
			"  %s [%d:%d-%d:%d]: %q\n",
			parseErr.Type,
			parseErr.StartRow+1,
			parseErr.StartCol+1,
			parseErr.EndRow+1,
			parseErr.EndCol+1,
			parseErr.Text,
		); err != nil {
			return fmt.Errorf("error writing parse summary: %w", err)
		}
	}
	return nil
}

func writeJSONParseSummary(w io.Writer, errors []jsonsanitize.ErrorNode) error {
	if len(errors) == 0 {
		if _, err := io.WriteString(w, "0 parse error(s)\n"); err != nil {
			return fmt.Errorf("error writing parse summary: %w", err)
		}
		return nil
	}

	if _, err := fmt.Fprintf(w, "%d parse error(s)\n", len(errors)); err != nil {
		return fmt.Errorf("error writing parse summary: %w", err)
	}
	for _, parseErr := range errors {
		if _, err := fmt.Fprintf(
			w,
			"  %s [%d:%d-%d:%d]: %q\n",
			parseErr.Type,
			parseErr.StartRow+1,
			parseErr.StartCol+1,
			parseErr.EndRow+1,
			parseErr.EndCol+1,
			parseErr.Text,
		); err != nil {
			return fmt.Errorf("error writing parse summary: %w", err)
		}
	}
	return nil
}

func normalizeFormat(format string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "yaml":
		return "yaml", nil
	case "json":
		return "json", nil
	default:
		return "", fmt.Errorf("unsupported format %q", format)
	}
}
