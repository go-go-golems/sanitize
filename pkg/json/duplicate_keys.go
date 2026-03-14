package jsonsanitize

import (
	stdjson "encoding/json"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

type duplicateKeyOccurrence struct {
	Key            string
	KeyText        string
	StartByte      uint
	EndByte        uint
	DuplicateIndex int
}

func collectDuplicateKeys(root *sitter.Node, content []byte) []duplicateKeyOccurrence {
	var duplicates []duplicateKeyOccurrence
	var visit func(*sitter.Node)
	visit = func(node *sitter.Node) {
		if node == nil {
			return
		}

		if node.Kind() == "object" {
			duplicates = append(duplicates, collectObjectDuplicates(node, content)...)
		}

		cursor := node.Walk()
		defer cursor.Close()
		children := node.NamedChildren(cursor)
		for i := range children {
			visit(&children[i])
		}
	}

	visit(root)
	return duplicates
}

func collectObjectDuplicates(node *sitter.Node, src []byte) []duplicateKeyOccurrence {
	cursor := node.Walk()
	defer cursor.Close()

	children := node.NamedChildren(cursor)
	seen := map[string]int{}
	var duplicates []duplicateKeyOccurrence

	for i := range children {
		child := &children[i]
		if child.Kind() != "pair" {
			continue
		}

		keyNode := child.ChildByFieldName("key")
		if keyNode == nil {
			continue
		}

		keyText := nodeText(src, keyNode)
		key := duplicateKeyIdentity(keyText)
		if key == "" {
			continue
		}

		seen[key]++
		if seen[key] > 1 {
			duplicates = append(duplicates, duplicateKeyOccurrence{
				Key:            key,
				KeyText:        keyText,
				StartByte:      keyNode.StartByte(),
				EndByte:        keyNode.EndByte(),
				DuplicateIndex: seen[key],
			})
		}
	}

	return duplicates
}

func duplicateKeyIdentity(keyText string) string {
	keyText = strings.TrimSpace(keyText)
	if keyText == "" {
		return ""
	}

	var unquoted string
	if err := stdjson.Unmarshal([]byte(keyText), &unquoted); err == nil {
		return unquoted
	}
	return keyText
}

func nodeText(src []byte, node *sitter.Node) string {
	startByte := node.StartByte()
	endByte := node.EndByte()
	if startByte >= endByte || endByte > uint(len(src)) {
		return ""
	}
	return string(src[startByte:endByte])
}
