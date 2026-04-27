package envfile

import (
	"testing"
)

func trimEntries() []Entry {
	return []Entry{
		{Key: "  APP_NAME  ", Value: "  myapp  "},
		{Key: "APP_SECRET", Value: `"s3cr3t"`},
		{Key: "PROD_HOST", Value: "localhost"},
		{Key: "PROD_PORT", Value: "'8080'"},
	}
}

func TestTrim_KeySpace(t *testing.T) {
	entries := trimEntries()
	res := Trim(entries, TrimOptions{TrimKeySpace: true})
	if res.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected trimmed key, got %q", res.Entries[0].Key)
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestTrim_ValueSpace(t *testing.T) {
	entries := trimEntries()
	res := Trim(entries, TrimOptions{TrimValueSpace: true})
	if res.Entries[0].Value != "myapp" {
		t.Errorf("expected trimmed value, got %q", res.Entries[0].Value)
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestTrim_Quotes_Double(t *testing.T) {
	entries := trimEntries()
	res := Trim(entries, TrimOptions{TrimQuotes: true})
	if res.Entries[1].Value != "s3cr3t" {
		t.Errorf("expected unquoted value, got %q", res.Entries[1].Value)
	}
}

func TestTrim_Quotes_Single(t *testing.T) {
	entries := trimEntries()
	res := Trim(entries, TrimOptions{TrimQuotes: true})
	if res.Entries[3].Value != "8080" {
		t.Errorf("expected unquoted value, got %q", res.Entries[3].Value)
	}
}

func TestTrim_KeyPrefix(t *testing.T) {
	entries := trimEntries()
	res := Trim(entries, TrimOptions{TrimPrefix: "PROD_"})
	if res.Entries[2].Key != "HOST" {
		t.Errorf("expected key without prefix, got %q", res.Entries[2].Key)
	}
	if res.Entries[0].Key == "HOST" {
		t.Error("non-prefixed key should be unchanged")
	}
	if res.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", res.Modified)
	}
}

func TestTrim_KeySuffix(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME_V2", Value: "foo"},
		{Key: "HOST", Value: "bar"},
	}
	res := Trim(entries, TrimOptions{TrimSuffix: "_V2"})
	if res.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected suffix stripped, got %q", res.Entries[0].Key)
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestTrim_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
	}
	res := Trim(entries, TrimOptions{})
	if res.Modified != 0 {
		t.Errorf("expected 0 modified, got %d", res.Modified)
	}
}

func TestFormatTrimSummary_Zero(t *testing.T) {
	r := TrimResult{Modified: 0}
	got := FormatTrimSummary(r)
	if got != "trim: no entries changed" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestFormatTrimSummary_One(t *testing.T) {
	r := TrimResult{Modified: 1}
	got := FormatTrimSummary(r)
	if got != "trim: 1 entry changed" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestFormatTrimSummary_Many(t *testing.T) {
	r := TrimResult{Modified: 5}
	got := FormatTrimSummary(r)
	if got != "trim: 5 entries changed" {
		t.Errorf("unexpected summary: %q", got)
	}
}
