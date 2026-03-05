package sanitize

// Example is a named YAML snippet with a description of the error it contains.
type Example struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	YAML        string `json:"yaml"`
}

// Examples is the built-in set of broken YAML snippets.
var Examples = []Example{
	{
		Name:        "Extra colon in plain scalar",
		Description: "A plain (unquoted) scalar value contains a colon, which the parser interprets as a nested key.",
		YAML:        "foobar: foba: sldkjf\nbaz: hello\n",
	},
	{
		Name:        "Tab indentation",
		Description: "YAML requires spaces for indentation; tabs cause a parse error.",
		YAML:        "server:\n\thost: localhost\n\tport: 8080\n",
	},
	{
		Name:        "Missing space after colon",
		Description: "Many parsers require a space after the colon in a mapping.",
		YAML:        "name:Alice\nage:30\ncity: Paris\n",
	},
	{
		Name:        "Unquoted hash in value",
		Description: "An unquoted # starts a comment, silently truncating the value.",
		YAML:        "password: abc#123\ntoken: xyz#456\n",
	},
	{
		Name:        "Bad list dash (no space)",
		Description: "A list item dash must be followed by a space.",
		YAML:        "items:\n  -apple\n  -banana\n  - cherry\n",
	},
	{
		Name:        "Trailing comma in flow mapping",
		Description: "A trailing comma inside a flow mapping is rejected by strict parsers.",
		YAML:        "point: { x: 1, y: 2, }\ncolor: { r: 255, g: 0, b: 128, }\n",
	},
	{
		Name:        "Duplicate keys",
		Description: "Two sibling keys share the same name; the second silently overwrites the first in most parsers.",
		YAML:        "config:\n  timeout: 30\n  timeout: 60\n  retries: 3\n",
	},
	{
		Name:        "Mixed indentation depth",
		Description: "Inconsistent indentation depth confuses the block structure parser.",
		YAML:        "a:\n  b: 1\n   c: 2\n  d: 3\n",
	},
	{
		Name:        "Multiple errors",
		Description: "A snippet combining several common mistakes at once.",
		YAML:        "service:\n\tname:my-app\n\tport: 8080\n\ttags:\n\t  -web\n\t  -backend\n\tenv: KEY: VALUE\n",
	},
	{
		Name:        "Valid YAML (no errors)",
		Description: "A well-formed YAML document — the tree should be error-free.",
		YAML:        "name: Alice\nage: 30\naddress:\n  street: 123 Main St\n  city: Springfield\nhobbies:\n  - reading\n  - cycling\n",
	},
}
