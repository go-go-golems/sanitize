package yamlsanitize

import (
	"fmt"

	yaml "github.com/tree-sitter-grammars/tree-sitter-yaml/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// newParser creates a configured tree-sitter parser for YAML.
func newParser() *sitter.Parser {
	parser := sitter.NewParser()
	lang := sitter.NewLanguage(yaml.Language())
	if err := parser.SetLanguage(lang); err != nil {
		panic("failed to set YAML language: " + err.Error())
	}
	return parser
}

// ParseTree returns the tree-sitter sexp representation of the YAML source,
// plus any ERROR/MISSING nodes found.
func ParseTree(src string) (string, []ErrorNode, error) {
	parser := newParser()
	defer parser.Close()

	content := []byte(src)
	tree := parser.Parse(content, nil)
	if tree == nil {
		return "", nil, fmt.Errorf("tree-sitter returned nil tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	errors := collectErrors(root, content)
	return root.ToSexp(), errors, nil
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
