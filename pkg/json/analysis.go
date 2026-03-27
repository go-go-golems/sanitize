package jsonsanitize

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

type documentAnalysis struct {
	Source           string
	TreeText         string
	ParseErrors      []ErrorNode
	StrictParseError error
	DuplicateKeys    []duplicateKeyOccurrence
	LineIndex        lineIndex
}

func analyzeDocument(src string) (documentAnalysis, error) {
	parser := newParser()
	defer parser.Close()

	content := []byte(src)
	tree := parser.Parse(content, nil)
	if tree == nil {
		return documentAnalysis{}, fmt.Errorf("tree-sitter returned nil tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	return analyzeParsedDocument(root, content), nil
}

func analyzeParsedDocument(root *sitter.Node, content []byte) documentAnalysis {
	return documentAnalysis{
		Source:           string(content),
		TreeText:         root.ToSexp(),
		ParseErrors:      collectErrors(root, content),
		StrictParseError: strictParseBytes(content),
		DuplicateKeys:    collectDuplicateKeys(root, content),
		LineIndex:        newLineIndex(string(content)),
	}
}
