package envfile

import (
	"testing"
)

func vEntries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestValidate_AllPresent(t *testing.T) {
	entries := vEntries("HOST", "localhost", "PORT", "8080")
	res := Validate(entries, []string{"HOST", "PORT"})
	if !res.OK() {
		t.Fatalf("expected OK, got errors: %s", res.String())
	}
}

func TestValidate_MissingKey(t *testing.T) {
	entries := vEntries("HOST", "localhost")
	res := Validate(entries, []string{"HOST", "PORT"})
	if res.OK() {
		t.Fatal("expected validation errors")
	}
	if len(res.Errors) != 1 || res.Errors[0].Key != "PORT" {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	entries := vEntries("HOST", "  ", "PORT", "8080")
	res := Validate(entries, []string{"HOST", "PORT"})
	if res.OK() {
		t.Fatal("expected validation error for empty value")
	}
	if res.Errors[0].Key != "HOST" {
		t.Fatalf("expected HOST error, got %v", res.Errors[0])
	}
}

func TestValidate_NoRequired(t *testing.T) {
	entries := vEntries("A", "1")
	res := Validate(entries, nil)
	if !res.OK() {
		t.Fatal("expected OK with no required keys")
	}
}

func TestValidationResult_String(t *testing.T) {
	res := ValidationResult{}
	if res.String() != "validation passed" {
		t.Fatalf("unexpected string: %s", res.String())
	}
	res.Errors = append(res.Errors, ValidationError{Key: "X", Message: "missing required key"})
	if res.String() == "validation passed" {
		t.Fatal("expected error string")
	}
}
