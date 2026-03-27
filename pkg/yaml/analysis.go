package yamlsanitize

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

type documentAnalysis struct {
	TreeText      string
	ParseErrors   []ErrorNode
	DuplicateKeys []duplicateKeyOccurrence
	LineIndex     lineIndex
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
		TreeText:      root.ToSexp(),
		ParseErrors:   collectErrors(root, content),
		DuplicateKeys: collectDuplicateKeys(root, content),
		LineIndex:     newLineIndex(string(content)),
	}
}
