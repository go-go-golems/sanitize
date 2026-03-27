package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	treejson "github.com/tree-sitter/tree-sitter-json/bindings/go"
)

type errorNode struct {
	Type      string `json:"type"`
	StartByte uint   `json:"start_byte"`
	EndByte   uint   `json:"end_byte"`
	StartRow  uint   `json:"start_row"`
	StartCol  uint   `json:"start_col"`
	EndRow    uint   `json:"end_row"`
	EndCol    uint   `json:"end_col"`
	Text      string `json:"text"`
}

type result struct {
	TreeText string      `json:"tree_text"`
	Errors   []errorNode `json:"errors"`
}

func main() {
	content, err := readInput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	res, err := parseJSON(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readInput() ([]byte, error) {
	if len(os.Args) > 1 {
		return os.ReadFile(os.Args[1])
	}
	return io.ReadAll(os.Stdin)
}

func parseJSON(content []byte) (result, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(treejson.Language())
	if err := parser.SetLanguage(lang); err != nil {
		return result{}, fmt.Errorf("set language: %w", err)
	}

	tree := parser.Parse(content, nil)
	if tree == nil {
		return result{}, fmt.Errorf("tree-sitter returned nil tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	return result{
		TreeText: root.ToSexp(),
		Errors:   collectErrors(root, content),
	}, nil
}

func collectErrors(node *sitter.Node, content []byte) []errorNode {
	var errs []errorNode
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
			errs = append(errs, errorNode{
				Type:      typ,
				StartByte: n.StartByte(),
				EndByte:   n.EndByte(),
				StartRow:  start.Row,
				StartCol:  start.Column,
				EndRow:    end.Row,
				EndCol:    end.Column,
				Text:      text,
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
