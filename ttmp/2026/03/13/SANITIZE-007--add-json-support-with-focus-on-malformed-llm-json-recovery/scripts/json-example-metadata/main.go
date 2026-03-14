package main

import (
	stdjson "encoding/json"
	"fmt"
	"os"

	"sanitize007scripts/internal/corpus"
)

func main() {
	enc := stdjson.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(corpus.LoadJSONExamples()); err != nil {
		fmt.Fprintf(os.Stderr, "encode metadata: %v\n", err)
		os.Exit(1)
	}
}
