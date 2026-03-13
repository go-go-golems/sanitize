package yamlsanitize

import (
	"fmt"
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

func findDuplicateKeys(src string) []duplicateKeyOccurrence {
	parser := newParser()
	defer parser.Close()

	content := []byte(src)
	tree := parser.Parse(content, nil)
	if tree == nil {
		return nil
	}
	defer tree.Close()

	var duplicates []duplicateKeyOccurrence
	var visit func(*sitter.Node)
	visit = func(node *sitter.Node) {
		if node == nil {
			return
		}

		if isMappingNode(node.Kind()) {
			duplicates = append(duplicates, collectMappingDuplicates(node, content)...)
		}

		cursor := node.Walk()
		defer cursor.Close()
		children := node.NamedChildren(cursor)
		for i := range children {
			visit(&children[i])
		}
	}

	visit(tree.RootNode())
	return duplicates
}

func collectMappingDuplicates(node *sitter.Node, src []byte) []duplicateKeyOccurrence {
	cursor := node.Walk()
	defer cursor.Close()

	children := node.NamedChildren(cursor)
	seen := map[string]int{}
	var duplicates []duplicateKeyOccurrence

	for i := range children {
		child := &children[i]
		if !isMappingPairNode(child.Kind()) {
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

func isMappingNode(kind string) bool {
	return kind == "block_mapping" || kind == "flow_mapping"
}

func isMappingPairNode(kind string) bool {
	return kind == "block_mapping_pair" || kind == "flow_pair"
}

func duplicateKeyIdentity(keyText string) string {
	keyText = strings.TrimSpace(keyText)
	if keyText == "" {
		return ""
	}
	if len(keyText) >= 2 {
		if (keyText[0] == '"' && keyText[len(keyText)-1] == '"') ||
			(keyText[0] == '\'' && keyText[len(keyText)-1] == '\'') {
			return keyText[1 : len(keyText)-1]
		}
	}
	return keyText
}

func duplicateKeyReplacement(keyText string, duplicateIndex int) string {
	suffix := fmt.Sprintf("_%d", duplicateIndex)
	if len(keyText) >= 2 {
		if (keyText[0] == '"' && keyText[len(keyText)-1] == '"') ||
			(keyText[0] == '\'' && keyText[len(keyText)-1] == '\'') {
			return keyText[:len(keyText)-1] + suffix + keyText[len(keyText)-1:]
		}
	}
	return keyText + suffix
}

func nodeText(src []byte, node *sitter.Node) string {
	startByte := node.StartByte()
	endByte := node.EndByte()
	if startByte >= endByte || endByte > uint(len(src)) {
		return ""
	}
	return string(src[startByte:endByte])
}

func lineIndexAtByte(src string, byteOffset uint) int {
	line := 0
	for i := 0; i < len(src) && uint(i) < byteOffset; i++ {
		if src[i] == '\n' {
			line++
		}
	}
	return line
}
