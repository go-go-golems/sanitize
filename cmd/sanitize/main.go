package main

import (
	"io"
	"os"

	sanitizecli "github.com/go-go-golems/sanitize/internal/cli"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return sanitizecli.Run(args, stdin, stdout, stderr)
}
