package utils

import (
	"encoding/json"
	"testing"
)

func TestNormalizeStringArrayFields_NumericAllowedPaths(t *testing.T) {
	input := `{"space_role":{"id":1,"allowed_paths":[325468791,343319186],"role":"test"}}`
	got := normalizeStringArrayFields([]byte(input))

	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	spaceRole, _ := result["space_role"].(map[string]interface{})
	paths, _ := spaceRole["allowed_paths"].([]interface{})

	if len(paths) != 2 {
		t.Fatalf("expected 2 allowed_paths, got %d", len(paths))
	}
	for i, expected := range []string{"325468791", "343319186"} {
		got, ok := paths[i].(string)
		if !ok {
			t.Errorf("allowed_paths[%d] is not a string: %T", i, paths[i])
			continue
		}
		if got != expected {
			t.Errorf("allowed_paths[%d]: expected %q, got %q", i, expected, got)
		}
	}
}

func TestNormalizeStringArrayFields_StringAllowedPaths(t *testing.T) {
	// String values must be left unchanged.
	input := `{"space_role":{"id":1,"allowed_paths":["325468791","343319186"]}}`
	got := normalizeStringArrayFields([]byte(input))

	// The byte output should be identical to the input (no re-encode needed
	// when nothing changed), but at minimum the values must still be strings.
	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	spaceRole, _ := result["space_role"].(map[string]interface{})
	paths, _ := spaceRole["allowed_paths"].([]interface{})

	if len(paths) != 2 {
		t.Fatalf("expected 2 allowed_paths, got %d", len(paths))
	}
	for i, expected := range []string{"325468791", "343319186"} {
		got, ok := paths[i].(string)
		if !ok {
			t.Errorf("allowed_paths[%d] is not a string: %T", i, paths[i])
			continue
		}
		if got != expected {
			t.Errorf("allowed_paths[%d]: expected %q, got %q", i, expected, got)
		}
	}
}

func TestNormalizeStringArrayFields_NumericResolvedAllowedPaths(t *testing.T) {
	input := `{"space_role":{"id":1,"resolved_allowed_paths":[111,222]}}`
	got := normalizeStringArrayFields([]byte(input))

	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	spaceRole, _ := result["space_role"].(map[string]interface{})
	paths, _ := spaceRole["resolved_allowed_paths"].([]interface{})

	if len(paths) != 2 {
		t.Fatalf("expected 2 resolved_allowed_paths, got %d", len(paths))
	}
	for i, expected := range []string{"111", "222"} {
		got, ok := paths[i].(string)
		if !ok {
			t.Errorf("resolved_allowed_paths[%d] is not a string: %T", i, paths[i])
			continue
		}
		if got != expected {
			t.Errorf("resolved_allowed_paths[%d]: expected %q, got %q", i, expected, got)
		}
	}
}

func TestNormalizeStringArrayFields_EmptyPaths(t *testing.T) {
	input := `{"space_role":{"id":1,"allowed_paths":[]}}`
	got := normalizeStringArrayFields([]byte(input))

	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	spaceRole, _ := result["space_role"].(map[string]interface{})
	paths, _ := spaceRole["allowed_paths"].([]interface{})
	if len(paths) != 0 {
		t.Fatalf("expected empty allowed_paths, got %d items", len(paths))
	}
}

func TestNormalizeStringArrayFields_InvalidJSON(t *testing.T) {
	input := []byte(`not json`)
	got := normalizeStringArrayFields(input)
	if string(got) != string(input) {
		t.Errorf("expected original bytes to be returned for invalid JSON")
	}
}

func TestNormalizeStringArrayFields_NoSpaceRoleKey(t *testing.T) {
	// Responses without a space_role key must be left untouched.
	input := `{"component":{"id":1}}`
	got := normalizeStringArrayFields([]byte(input))
	if string(got) != input {
		t.Errorf("expected no change, got %s", got)
	}
}
