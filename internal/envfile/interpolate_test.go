package envfile

import (
	"testing"
)

func interpEntries() []Entry {
	return []Entry{
		{Key: "BASE_URL", Value: "https://example.com"},
		{Key: "API_URL", Value: "${BASE_URL}/api"},
		{Key: "CALLBACK", Value: "$BASE_URL/callback"},
		{Key: "PLAIN", Value: "no-refs-here"},
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	res, err := Interpolate(interpEntries(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://example.com/api"
	if res.Entries[1].Value != want {
		t.Errorf("got %q, want %q", res.Entries[1].Value, want)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	res, err := Interpolate(interpEntries(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://example.com/callback"
	if res.Entries[2].Value != want {
		t.Errorf("got %q, want %q", res.Entries[2].Value, want)
	}
}

func TestInterpolate_ExpandedCount(t *testing.T) {
	res, err := Interpolate(interpEntries(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expanded != 2 {
		t.Errorf("expanded count: got %d, want 2", res.Expanded)
	}
}

func TestInterpolate_UnresolvedNonStrict(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "${MISSING_VAR}"},
	}
	res, err := Interpolate(entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "MISSING_VAR" {
		t.Errorf("expected MISSING_VAR in unresolved, got %v", res.Unresolved)
	}
	if res.Entries[0].Value != "${MISSING_VAR}" {
		t.Errorf("value should be unchanged, got %q", res.Entries[0].Value)
	}
}

func TestInterpolate_UnresolvedStrict(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "${MISSING_VAR}"},
	}
	_, err := Interpolate(entries, true)
	if err == nil {
		t.Error("expected error for unresolved reference in strict mode")
	}
}

func TestInterpolate_PlainValue(t *testing.T) {
	res, err := Interpolate(interpEntries(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[3].Value != "no-refs-here" {
		t.Errorf("plain value mutated: %q", res.Entries[3].Value)
	}
}
