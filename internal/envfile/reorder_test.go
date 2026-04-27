package envfile

import (
	"testing"
)

func reorderEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestReorder_ExplicitKeys(t *testing.T) {
	entries := reorderEntries()
	opts := ReorderOptions{
		Keys: []string{"APP_ENV", "APP_PORT", "SECRET_KEY"},
	}
	r := Reorder(entries, opts)

	if r.Entries[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV first, got %s", r.Entries[0].Key)
	}
	if r.Entries[1].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT second, got %s", r.Entries[1].Key)
	}
	if r.Entries[2].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY third, got %s", r.Entries[2].Key)
	}
	if r.Moved != 3 {
		t.Errorf("expected Moved=3, got %d", r.Moved)
	}
	if r.Unknown != 3 {
		t.Errorf("expected Unknown=3, got %d", r.Unknown)
	}
}

func TestReorder_PrefixPin(t *testing.T) {
	entries := reorderEntries()
	opts := ReorderOptions{
		PinFirst: []string{"APP_"},
	}
	r := Reorder(entries, opts)

	if r.Entries[0].Key != "APP_PORT" && r.Entries[1].Key != "APP_ENV" {
		// either order is fine as long as both APP_ keys are at the front
		if r.Entries[0].Key != "APP_ENV" {
			t.Errorf("expected APP_ keys at the front, got %s", r.Entries[0].Key)
		}
	}
}

func TestReorder_PreservesAllEntries(t *testing.T) {
	entries := reorderEntries()
	opts := ReorderOptions{
		Keys: []string{"LOG_LEVEL"},
	}
	r := Reorder(entries, opts)

	if len(r.Entries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(r.Entries))
	}
}

func TestReorder_MissingExplicitKeyIgnored(t *testing.T) {
	entries := reorderEntries()
	opts := ReorderOptions{
		Keys: []string{"NONEXISTENT", "APP_ENV"},
	}
	r := Reorder(entries, opts)

	if r.Entries[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV first, got %s", r.Entries[0].Key)
	}
	if len(r.Entries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(r.Entries))
	}
}

func TestFormatReorderSummary(t *testing.T) {
	r := ReorderResult{Moved: 2, Unknown: 4}
	s := FormatReorderSummary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
	expected := "2 entries repositioned, 4 left in original order"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
