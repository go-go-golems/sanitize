package main

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	treejson "github.com/tree-sitter/tree-sitter-json/bindings/go"
)

type testCase struct {
	Number       string
	Title        string
	Language     string
	Snippet      string
	ExpectedNote string
}

type parseErr struct {
	Type     string
	StartRow uint
	StartCol uint
	EndRow   uint
	EndCol   uint
	Text     string
}

func main() {
	input := defaultInput()
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	content, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input: %v\n", err)
		os.Exit(1)
	}

	cases := parseCases(string(content))
	fmt.Println("---")
	fmt.Println("Title: JSON parse error replication matrix")
	fmt.Println("Ticket: SANITIZE-007")
	fmt.Println("Status: active")
	fmt.Println("Topics:")
	fmt.Println("    - json")
	fmt.Println("    - linting")
	fmt.Println("    - api-design")
	fmt.Println("DocType: reference")
	fmt.Println("Intent: long-term")
	fmt.Println("Owners: []")
	fmt.Println("RelatedFiles: []")
	fmt.Println("ExternalSources: []")
	fmt.Println("Summary: Replication matrix comparing strict encoding/json failures with tree-sitter JSON error nodes for the imported malformed LLM JSON cases.")
	fmt.Println("WhatFor: Show which malformed JSON cases tree-sitter can localize structurally and where heuristics will still be required.")
	fmt.Println("WhenToUse: Use when designing JSON parse-aware lint rules and fix heuristics.")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# JSON Parse Error Replication Matrix")
	fmt.Println()
	fmt.Printf("Source: `%s`\n\n", input)

	for _, tc := range cases {
		fmt.Printf("## Case %s: %s\n\n", tc.Number, tc.Title)
		if tc.Snippet == "" {
			fmt.Println("- Snippet: none in source note")
			if tc.ExpectedNote != "" {
				fmt.Printf("- Source note: %s\n", tc.ExpectedNote)
			}
			fmt.Println()
			continue
		}

		strictValid, strictErr := strictParse(tc.Snippet)
		tree, errs, treeErr := parseTree(tc.Snippet)

		fmt.Printf("- Language tag: `%s`\n", tc.Language)
		fmt.Printf("- Strict `encoding/json` valid: `%t`\n", strictValid)
		if strictErr != "" {
			fmt.Printf("- Strict parse error: `%s`\n", strictErr)
		}
		if treeErr != nil {
			fmt.Printf("- Tree-sitter parser setup error: `%v`\n", treeErr)
		} else {
			fmt.Printf("- Tree-sitter error node count: `%d`\n", len(errs))
			if len(errs) > 0 {
				first := errs[0]
				fmt.Printf("- First tree-sitter error: `%s` rows %d:%d to %d:%d text=%q\n", first.Type, first.StartRow, first.StartCol, first.EndRow, first.EndCol, first.Text)
			}
		}
		if tc.ExpectedNote != "" {
			fmt.Printf("- Source note: %s\n", tc.ExpectedNote)
		}
		fmt.Printf("- Heuristic hits: `%s`\n", strings.Join(detectHeuristics(tc.Snippet), ", "))
		fmt.Println()

		snippetFence := renderFence(tc.Snippet)
		fmt.Printf("%stext\n", snippetFence)
		fmt.Println(tc.Snippet)
		fmt.Println(snippetFence)
		fmt.Println()

		fmt.Println("```text")
		fmt.Println(tree)
		fmt.Println("```")
		fmt.Println()
	}
}

func defaultInput() string {
	return filepath.Join("..", "sources", "local", "03-json-parse-errors-import.md")
}

func parseCases(content string) []testCase {
	headerRe := regexp.MustCompile(`^(\d+)\.\s+\*\*(.+?)\*\*`)
	lines := strings.Split(content, "\n")
	var cases []testCase
	var current *testCase

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if m := headerRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				cases = append(cases, *current)
			}
			current = &testCase{Number: m[1], Title: m[2]}
			continue
		}
		if current == nil {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Error:") {
			current.ExpectedNote = strings.TrimPrefix(trimmed, "Error: ")
			continue
		}

		if fence, lang, ok := fenceParts(trimmed); ok {
			var block []string
			for i = i + 1; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == fence {
					break
				}
				block = append(block, lines[i])
			}
			current.Language = lang
			current.Snippet = strings.Join(block, "\n")
		}
	}
	if current != nil {
		cases = append(cases, *current)
	}
	return cases
}

func fenceParts(trimmed string) (string, string, bool) {
	if !strings.HasPrefix(trimmed, "```") {
		return "", "", false
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == '`' {
		n++
	}
	if n < 3 {
		return "", "", false
	}
	return strings.Repeat("`", n), strings.TrimSpace(trimmed[n:]), true
}

func strictParse(snippet string) (bool, string) {
	dec := stdjson.NewDecoder(strings.NewReader(snippet))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return false, err.Error()
	}
	var extra any
	if err := dec.Decode(&extra); err == nil {
		return false, "multiple top-level JSON values"
	}
	return true, ""
}

func parseTree(snippet string) (string, []parseErr, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(treejson.Language())
	if err := parser.SetLanguage(lang); err != nil {
		return "", nil, err
	}

	content := []byte(snippet)
	tree := parser.Parse(content, nil)
	if tree == nil {
		return "", nil, fmt.Errorf("tree-sitter returned nil tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	return root.ToSexp(), collectErrors(root, content), nil
}

func collectErrors(node *sitter.Node, content []byte) []parseErr {
	var errs []parseErr
	var walk func(*sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.IsError() || n.IsMissing() {
			typ := "ERROR"
			if n.IsMissing() {
				typ = "MISSING"
			}
			start := n.StartPosition()
			end := n.EndPosition()
			text := ""
			if sb, eb := n.StartByte(), n.EndByte(); sb < eb && eb <= uint(len(content)) {
				text = string(content[sb:eb])
				text = strings.ReplaceAll(text, "\n", `\n`)
			}
			errs = append(errs, parseErr{
				Type:     typ,
				StartRow: start.Row,
				StartCol: start.Column,
				EndRow:   end.Row,
				EndCol:   end.Column,
				Text:     text,
			})
		}
		cursor := n.Walk()
		defer cursor.Close()
		children := n.Children(cursor)
		for i := range children {
			walk(&children[i])
		}
	}
	walk(node)
	return errs
}

func detectHeuristics(snippet string) []string {
	hits := []string{}
	add := func(name string, ok bool) {
		if ok {
			hits = append(hits, name)
		}
	}

	add("markdown_fence", strings.Contains(snippet, "```"))
	add("comment", strings.Contains(snippet, "//") || strings.Contains(snippet, "/*"))
	add("single_quotes", strings.Contains(snippet, "'"))
	add("python_literals", strings.Contains(snippet, "True") || strings.Contains(snippet, "False") || strings.Contains(snippet, "None"))
	add("trailing_comma", regexp.MustCompile(`,\s*[}\]]`).MatchString(snippet))
	add("ellipsis_or_placeholder", strings.Contains(snippet, "...") || strings.Contains(snippet, "<") || strings.Contains(snippet, "[insert"))
	add("duplicate_comma", strings.Contains(snippet, ",,"))
	add("prose_wrapped", strings.Contains(snippet, "Here is your JSON:") || strings.Contains(snippet, "Thanks!"))
	add("multiple_top_level", bytes.Count([]byte(snippet), []byte("}{")) > 0 || strings.Contains(snippet, "} {"))

	if len(hits) == 0 {
		hits = append(hits, "none")
	}
	return hits
}

func renderFence(snippet string) string {
	if strings.Contains(snippet, "```") {
		return "````"
	}
	return "```"
}
