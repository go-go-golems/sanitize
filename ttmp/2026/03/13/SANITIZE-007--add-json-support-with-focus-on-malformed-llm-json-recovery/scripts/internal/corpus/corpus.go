package corpus

import (
	"sort"
	"strings"

	examplespkg "github.com/go-go-golems/sanitize/examples"
	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
)

type Example struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Input       string `json:"input"`
	Category    string `json:"category,omitempty"`
	Source      string `json:"source,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
}

func LoadJSONExamples() []Example {
	ret := make([]Example, 0, len(jsonsanitize.Examples)+len(examplespkg.LoadJSONExamples()))

	for _, ex := range jsonsanitize.Examples {
		ret = append(ret, Example{
			Name:        ex.Name,
			Description: ex.Description,
			Input:       ex.JSON,
			Category:    "demo",
			Source:      "builtin",
			Pattern:     patternForExample(ex.Filename, ex.Name),
		})
	}

	for _, ex := range examplespkg.LoadJSONExamples() {
		ret = append(ret, Example{
			Name:        ex.Name,
			Description: ex.Description,
			Input:       ex.JSON,
			Category:    ex.Category,
			Source:      ex.Source,
			Filename:    ex.Filename,
			Pattern:     patternForExample(ex.Filename, ex.Name),
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		left := ret[i].Source + "\x00" + ret[i].Filename + "\x00" + ret[i].Name
		right := ret[j].Source + "\x00" + ret[j].Filename + "\x00" + ret[j].Name
		return left < right
	})
	return ret
}

func patternForExample(filename, name string) string {
	switch filename {
	case "10-trailing-comma-object.json":
		return "trailing_comma"
	case "11-single-quotes.json":
		return "single_quotes"
	case "12-unquoted-keys.json":
		return "unquoted_keys"
	case "13-missing-comma.json":
		return "missing_comma"
	case "14-missing-colon.json":
		return "missing_colon"
	case "15-leading-prose-wrapper.json":
		return "leading_or_trailing_prose"
	case "16-markdown-fence-wrapper.json":
		return "markdown_fence_wrapper"
	case "17-python-literals.json":
		return "python_literals"
	case "18-duplicate-comma-array.json":
		return "duplicate_comma"
	case "19-comments.json":
		return "comment"
	case "20-missing-closing-delimiter.json":
		return "missing_closing_delimiter"
	case "21-multiple-top-level-objects.json":
		return "multiple_top_level_values"
	case "22-unterminated-string.json":
		return "unterminated_string"
	case "23-duplicate-key.json":
		return "duplicate_key"
	case "24-llm-wrapper-multi-step.json":
		return "wrapper_python_trailing_comma_combo"
	case "25-llm-commentary-comments-and-duplicate-comma.json":
		return "comment_python_duplicate_comma_combo"
	default:
		switch strings.ToLower(name) {
		case "trailing comma":
			return "trailing_comma"
		case "single quotes":
			return "single_quotes"
		case "unquoted keys":
			return "unquoted_keys"
		case "python literals":
			return "python_literals"
		case "markdown fence wrapper":
			return "markdown_fence_wrapper"
		case "leading prose":
			return "leading_or_trailing_prose"
		case "multiple top-level values":
			return "multiple_top_level_values"
		case "valid json":
			return "valid"
		default:
			return strings.ReplaceAll(strings.ToLower(name), " ", "_")
		}
	}
}
