package sanitizecli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	glazecmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/sanitize/internal/server"
	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

type fixSettings struct {
	Input         string `glazed:"input"`
	JSON          bool   `glazed:"json"`
	TabWidth      int    `glazed:"tab-width"`
	MaxIterations int    `glazed:"max-iterations"`
}

type lintSettings struct {
	Input string `glazed:"input"`
	JSON  bool   `glazed:"json"`
}

type parseSettings struct {
	Input string `glazed:"input"`
	JSON  bool   `glazed:"json"`
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
				fields.New("json", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Output the full sanitize result as JSON")),
				fields.New("tab-width", fields.TypeInteger, fields.WithDefault(2), fields.WithHelp("Spaces per tab for the tab_indent fixer")),
				fields.New("max-iterations", fields.TypeInteger, fields.WithDefault(10), fields.WithHelp("Maximum number of fix iterations")),
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

	input, err := readInput(settings.Input, c.streams.Stdin)
	if err != nil {
		return newExitError(1, fmt.Errorf("error reading input: %w", err))
	}

	result := yamlsanitize.Sanitize(string(input),
		yamlsanitize.WithTabWidth(settings.TabWidth),
		yamlsanitize.WithMaxIterations(settings.MaxIterations),
	)

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
		if err := writeFixSummary(c.streams.Stderr, result.Fixes); err != nil {
			return newExitError(1, err)
		}
	}

	if !result.ParseClean || !result.LintClean {
		return newExitError(1, nil)
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
				fields.New("json", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Output lint issues as JSON")),
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

	input, err := readInput(settings.Input, c.streams.Stdin)
	if err != nil {
		return newExitError(1, fmt.Errorf("error reading input: %w", err))
	}

	issues := yamlsanitize.Lint(string(input))
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

	input, err := readInput(settings.Input, c.streams.Stdin)
	if err != nil {
		return newExitError(1, fmt.Errorf("error reading input: %w", err))
	}

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
		if err := writeParseSummary(c.streams.Stderr, result.Errors); err != nil {
			return newExitError(1, err)
		}
	}

	if len(result.Errors) > 0 {
		return newExitError(1, nil)
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

func readInput(path string, stdin io.Reader) ([]byte, error) {
	if path != "" {
		return os.ReadFile(path)
	}
	return io.ReadAll(stdin)
}

func writeFixSummary(w io.Writer, fixes []yamlsanitize.Fix) error {
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

func writeParseSummary(w io.Writer, errors []yamlsanitize.ErrorNode) error {
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
