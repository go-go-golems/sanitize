package examples

import "testing"

func TestLoadFileExamples(t *testing.T) {
	examples := LoadFileExamples()
	if len(examples) == 0 {
		t.Fatal("expected YAML examples")
	}
	if examples[0].Filename == "" {
		t.Fatalf("expected filenames to be populated, got %+v", examples[0])
	}
}

func TestLoadJSONExamples(t *testing.T) {
	examples := LoadJSONExamples()
	if len(examples) == 0 {
		t.Fatal("expected JSON examples")
	}
	if examples[0].Filename == "" {
		t.Fatalf("expected filenames to be populated, got %+v", examples[0])
	}
	if examples[0].JSON == "" {
		t.Fatalf("expected example payload, got %+v", examples[0])
	}
}
