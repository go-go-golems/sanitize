package jsonsanitize

import "testing"

func TestSanitizeWithOptions_StripsFenceAndNormalizesLiteralAndComma(t *testing.T) {
	src := "```json\n{\"ok\": True,}\n```"

	result, err := SanitizeWithOptions(src)
	if err != nil {
		t.Fatalf("SanitizeWithOptions: %v", err)
	}

	if result.Sanitized != "{\"ok\": true}\n" {
		t.Fatalf("unexpected sanitized output:\n%s", result.Sanitized)
	}
	if !result.ParseClean || !result.StrictParseClean {
		t.Fatalf("expected parse clean result, got %+v", result)
	}
	if len(result.Fixes) == 0 {
		t.Fatal("expected fixes to be recorded")
	}
}

func TestSanitizeWithOptions_StripsCommentsAndDuplicateComma(t *testing.T) {
	src := "{\"items\":[1,,2], // hi\n\"ok\": true\n}\n"

	result, err := SanitizeWithOptions(src)
	if err != nil {
		t.Fatalf("SanitizeWithOptions: %v", err)
	}

	if result.Sanitized != "{\"items\":[1,2], \n\"ok\": true\n}\n" {
		t.Fatalf("unexpected sanitized output:\n%q", result.Sanitized)
	}
	if len(result.Fixes) < 2 {
		t.Fatalf("expected multiple fixes, got %+v", result.Fixes)
	}
}

func TestSanitizeWithOptions_ExtractsJSONFromProse(t *testing.T) {
	src := "Here is your JSON:\n{\"name\":\"Alice\"}\nThanks!"

	result, err := SanitizeWithOptions(src)
	if err != nil {
		t.Fatalf("SanitizeWithOptions: %v", err)
	}

	if result.Sanitized != "{\"name\":\"Alice\"}\n" {
		t.Fatalf("unexpected sanitized output:\n%q", result.Sanitized)
	}
	if !result.StrictParseClean {
		t.Fatalf("expected strict parse clean result, got %+v", result)
	}
}

func TestSanitizeWithOptions_DisabledRuleSkipsFix(t *testing.T) {
	src := "{\"ok\": True}\n"

	result, err := SanitizeWithOptions(src, WithDisabledRules("python_literals"))
	if err != nil {
		t.Fatalf("SanitizeWithOptions: %v", err)
	}

	if result.Sanitized != src {
		t.Fatalf("expected python literal to remain unchanged, got %q", result.Sanitized)
	}
}
