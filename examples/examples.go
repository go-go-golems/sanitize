// Package examples embeds the YAML example corpus and provides a loader.
package examples

import (
	"embed"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

//go:embed yaml/*.yaml
var yamlFiles embed.FS

// LoadFileExamples reads all embedded YAML files and returns them as Examples.
// Category is derived from the filename prefix per the corpus convention:
//
//	00-09 → "valid"
//	10-19 → "single-error"
//	20-29 → "mixed"
func LoadFileExamples() []yamlsanitize.Example {
	entries, err := yamlFiles.ReadDir("yaml")
	if err != nil {
		return nil
	}

	var examples []yamlsanitize.Example
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}

		data, err := yamlFiles.ReadFile(filepath.Join("yaml", e.Name()))
		if err != nil {
			continue
		}

		name, category := parseFilename(e.Name())
		examples = append(examples, yamlsanitize.Example{
			Name:     name,
			YAML:     string(data),
			Category: category,
			Source:   "file",
			Filename: e.Name(),
		})
	}

	sort.Slice(examples, func(i, j int) bool {
		return examples[i].Filename < examples[j].Filename
	})
	return examples
}

// parseFilename converts "11-missing-space-after-colon.yaml" into
// name="Missing space after colon", category="single-error".
func parseFilename(filename string) (string, string) {
	base := strings.TrimSuffix(filename, ".yaml")
	name := ""
	category := ""

	// Split on first dash: "11-missing-space-after-colon" → "11", "missing-space-after-colon"
	parts := strings.SplitN(base, "-", 2)
	if len(parts) == 2 {
		name = strings.ReplaceAll(parts[1], "-", " ")
		// Capitalize first letter
		if len(name) > 0 {
			name = strings.ToUpper(name[:1]) + name[1:]
		}

		num, err := strconv.Atoi(parts[0])
		if err == nil {
			switch {
			case num < 10:
				category = "valid"
			case num < 20:
				category = "single-error"
			default:
				category = "mixed"
			}
		}
	} else {
		name = base
	}
	return name, category
}
