package jsonsanitize

// Examples is the built-in set of broken or representative JSON snippets.
var Examples = []Example{
	{
		Name:        "Trailing comma",
		Description: "A trailing comma appears after the last object member.",
		JSON:        "{\"name\":\"Alice\",}\n",
	},
	{
		Name:        "Single quotes",
		Description: "Single quotes are used where JSON requires double quotes.",
		JSON:        "{'name': 'Alice'}\n",
	},
	{
		Name:        "Unquoted keys",
		Description: "Object keys are not quoted JSON strings.",
		JSON:        "{name: \"Alice\", age: 30}\n",
	},
	{
		Name:        "Python literals",
		Description: "Python-style literals appear in place of JSON literals.",
		JSON:        "{\"ok\": True, \"value\": None}\n",
	},
	{
		Name:        "Markdown fence wrapper",
		Description: "The JSON payload is wrapped in a Markdown code fence.",
		JSON:        "```json\n{\"name\":\"Alice\"}\n```\n",
	},
	{
		Name:        "Leading prose",
		Description: "Helpful prose appears before the JSON payload.",
		JSON:        "Here is your JSON:\n{\"name\":\"Alice\"}\n",
	},
	{
		Name:        "Multiple top-level values",
		Description: "Two top-level JSON objects appear back to back.",
		JSON:        "{\"a\":1} {\"b\":2}\n",
	},
	{
		Name:        "Valid JSON",
		Description: "A well-formed JSON document with no parse errors.",
		JSON:        "{\"name\":\"Alice\",\"age\":30,\"tags\":[\"ops\",\"infra\"]}\n",
	},
}
