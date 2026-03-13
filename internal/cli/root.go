package sanitizecli

import (
	"errors"
	"fmt"
	"io"

	glazecli "github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/spf13/cobra"
)

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type exitError struct {
	Code int
	Err  error
}

func (e *exitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("command exited with status %d", e.Code)
}

func (e *exitError) Unwrap() error {
	return e.Err
}

func newExitError(code int, err error) error {
	return &exitError{Code: code, Err: err}
}

func defaultSections() ([]schema.Section, error) {
	commandSettingsSection, err := glazecli.NewCommandSettingsSection()
	if err != nil {
		return nil, err
	}
	return []schema.Section{commandSettingsSection}, nil
}

func buildCobra(command cmds.Command) (*cobra.Command, error) {
	cobraCmd := glazecli.NewCobraCommandFromCommandDescription(command.Description())
	parser, err := glazecli.NewCobraParserFromSections(command.Description().Schema, &glazecli.CobraParserConfig{
		ShortHelpSections:          []string{schema.DefaultSlug},
		MiddlewaresFunc:            glazecli.CobraCommandDefaultMiddlewares,
		SkipCommandSettingsSection: true,
	})
	if err != nil {
		return nil, err
	}
	if err := parser.AddToCobraCommand(cobraCmd); err != nil {
		return nil, err
	}

	cobraCmd.RunE = func(cmd *cobra.Command, args []string) error {
		parsedValues, err := parser.Parse(cmd, args)
		if err != nil {
			return err
		}

		if bareCmd, ok := command.(cmds.BareCommand); ok {
			return bareCmd.Run(cmd.Context(), parsedValues)
		}
		if writerCmd, ok := command.(cmds.WriterCommand); ok {
			return writerCmd.RunIntoWriter(cmd.Context(), parsedValues, cmd.OutOrStdout())
		}
		return fmt.Errorf("command %q does not implement a supported run interface", command.Description().Name)
	}

	return cobraCmd, nil
}

func NewRootCommand(streams Streams) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:               "sanitize",
		Short:             "YAML linter and heuristic fixer",
		SilenceUsage:      true,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return logging.InitLoggerFromCobra(cmd)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	rootCmd.SetIn(streams.Stdin)
	rootCmd.SetOut(streams.Stdout)
	rootCmd.SetErr(streams.Stderr)

	if err := logging.AddLoggingSectionToRootCommand(rootCmd, "sanitize"); err != nil {
		return nil, err
	}

	constructors := []func(Streams) (cmds.Command, error){
		newFixCommand,
		newLintCommand,
		newRulesCommand,
		newParseCommand,
		newServeCommand,
	}

	for _, constructor := range constructors {
		command, err := constructor(streams)
		if err != nil {
			return nil, err
		}
		cobraCmd, err := buildCobra(command)
		if err != nil {
			return nil, err
		}
		rootCmd.AddCommand(cobraCmd)
	}

	return rootCmd, nil
}

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	rootCmd, err := NewRootCommand(Streams{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}

	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		var exitErr *exitError
		if errors.As(err, &exitErr) {
			if exitErr.Err != nil {
				_, _ = fmt.Fprintln(stderr, exitErr.Err)
			}
			return exitErr.Code
		}
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	return 0
}
