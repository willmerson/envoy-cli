package envfile

import (
	"testing"
)

func flatEntries() []Entry {
	return []Entry{
		{Key: "APP__DB__HOST", Value: "localhost"},
		{Key: "APP__DB__PORT", Value: "5432"},
		{Key: "APP__SECRET", Value: ""},
		{Key: "APP__API_KEY", Value: "abc123"},
		{Key: "OTHER_KEY", Value: "value"},
	}
}

func TestFlatten_CollapsesDoubleSeparator(t *testing.T) {
	entries := flatEntries()
	r := Flatten(entries, FlattenOptions{})

	if r.Entries[0].Key != "APP_DB_HOST" {
		t.Errorf("expected APP_DB_HOST, got %s", r.Entries[0].Key)
	}
	if r.Entries[1].Key != "APP_DB_PORT" {
		t.Errorf("expected APP_DB_PORT, got %s", r.Entries[1].Key)
	}
}

func TestFlatten_DropEmpty(t *testing.T) {
	entries := flatEntries()
	r := Flatten(entries, FlattenOptions{DropEmpty: true})

	for _, e := range r.Entries {
		if e.Value == "" {
			t.Errorf("expected empty entry to be dropped, got key %s", e.Key)
		}
	}
	if r.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", r.Dropped)
	}
}

func TestFlatten_StripPrefix(t *testing.T) {
	entries := flatEntries()
	r := Flatten(entries, FlattenOptions{StripPrefix: "APP"})

	for _, e := range r.Entries {
		if e.Key == "OTHER_KEY" {
			continue // not prefixed
		}
		if len(e.Key) > 3 && e.Key[:3] == "APP" {
			t.Errorf("prefix not stripped from key %s", e.Key)
		}
	}
	if r.Renamed == 0 {
		t.Error("expected at least one rename from prefix stripping")
	}
}

func TestFlatten_UppercaseKeys(t *testing.T) {
	entries := []Entry{
		{Key: "app_host", Value: "localhost"},
		{Key: "app_port", Value: "8080"},
	}
	r := Flatten(entries, FlattenOptions{UppercaseKeys: true})

	for _, e := range r.Entries {
		if e.Key != "APP_HOST" && e.Key != "APP_PORT" {
			t.Errorf("expected uppercase key, got %s", e.Key)
		}
	}
	if r.Renamed != 2 {
		t.Errorf("expected 2 renamed, got %d", r.Renamed)
	}
}

func TestFlatten_DeduplicatesAfterNormalise(t *testing.T) {
	entries := []Entry{
		{Key: "APP__KEY", Value: "first"},
		{Key: "APP_KEY", Value: "second"},
	}
	r := Flatten(entries, FlattenOptions{})

	if len(r.Entries) != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", len(r.Entries))
	}
	if r.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", r.Dropped)
	}
}

func TestFormatFlattenSummary(t *testing.T) {
	r := FlattenResult{Entries: make([]Entry, 3), Renamed: 2, Dropped: 1}
	summary := FormatFlattenSummary(r)
	expected := "3 entries kept, 2 renamed, 1 dropped"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
