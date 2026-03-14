package jsonsanitize

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"io"

	sitter "github.com/tree-sitter/go-tree-sitter"
	treejson "github.com/tree-sitter/tree-sitter-json/bindings/go"
)

// newParser creates a configured tree-sitter parser for JSON.
func newParser() *sitter.Parser {
	parser := sitter.NewParser()
	lang := sitter.NewLanguage(treejson.Language())
	if err := parser.SetLanguage(lang); err != nil {
		panic("failed to set JSON language: " + err.Error())
	}
	return parser
}

// ParseTree returns the tree-sitter sexp representation of the JSON source,
// plus any ERROR/MISSING nodes found.
func ParseTree(src string) (string, []ErrorNode, error) {
	analysis, err := analyzeDocument(src)
	if err != nil {
		return "", nil, err
	}
	return analysis.TreeText, analysis.ParseErrors, nil
}

// StrictParse validates that src is strict JSON with exactly one top-level value.
func StrictParse(src string) error {
	return strictParseBytes([]byte(src))
}

func strictParseBytes(content []byte) error {
	dec := stdjson.NewDecoder(bytes.NewReader(content))
	dec.UseNumber()

	var value any
	if err := dec.Decode(&value); err != nil {
		return err
	}

	var extra any
	switch err := dec.Decode(&extra); err {
	case io.EOF:
		return nil
	case nil:
		return fmt.Errorf("multiple top-level JSON values")
	default:
		return err
	}
}

// collectErrors walks the tree and gathers all ERROR and MISSING nodes.
func collectErrors(node *sitter.Node, src []byte) []ErrorNode {
	var errs []ErrorNode
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.IsError() || n.IsMissing() {
			typ := "ERROR"
			if n.IsMissing() {
				typ = "MISSING"
			}
			sp := n.StartPosition()
			ep := n.EndPosition()
			text := ""
			sb := n.StartByte()
			eb := n.EndByte()
			if sb < eb && eb <= uint(len(src)) {
				text = string(src[sb:eb])
			}
			errs = append(errs, ErrorNode{
				Type:      typ,
				StartByte: sb,
				EndByte:   eb,
				StartRow:  sp.Row,
				StartCol:  sp.Column,
				EndRow:    ep.Row,
				EndCol:    ep.Column,
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
