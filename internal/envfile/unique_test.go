package envfile

import (
	"strings"
	"testing"
)

func uniqueEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "CACHE_HOST", Value: "localhost"},
		{Key: "APP_NAME", Value: "duplicate-key"},
		{Key: "LOG_LEVEL", Value: "info"},
		{Key: "DEBUG", Value: "INFO"},
	}
}

func TestUnique_ByKey(t *testing.T) {
	entries := uniqueEntries()
	result := Unique(entries, UniqueOptions{ByKey: true})

	if result.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", result.Removed)
	}
	for _, e := range result.Entries {
		if e.Key == "APP_NAME" && e.Value == "duplicate-key" {
			t.Error("expected duplicate APP_NAME to be removed")
		}
	}
}

func TestUnique_ByValue_CaseSensitive(t *testing.T) {
	entries := uniqueEntries()
	result := Unique(entries, UniqueOptions{ByValue: true, CaseSensitive: true})

	// "localhost" appears twice; "info" and "INFO" are different when case-sensitive
	if result.Removed != 1 {
		t.Fatalf("expected 1 removed (localhost dup), got %d", result.Removed)
	}
}

func TestUnique_ByValue_CaseInsensitive(t *testing.T) {
	entries := uniqueEntries()
	result := Unique(entries, UniqueOptions{ByValue: true, CaseSensitive: false})

	// "localhost" x2 and "info"/"INFO" x2 => 2 removed
	if result.Removed != 2 {
		t.Fatalf("expected 2 removed, got %d", result.Removed)
	}
}

func TestUnique_NoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}
	result := Unique(entries, UniqueOptions{ByKey: true, ByValue: true})

	if result.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", result.Removed)
	}
	if len(result.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result.Entries))
	}
}

func TestUnique_Empty(t *testing.T) {
	result := Unique([]Entry{}, UniqueOptions{ByKey: true})
	if result.Removed != 0 || len(result.Entries) != 0 {
		t.Error("expected empty result for empty input")
	}
}

func TestFormatUniqueSummary_WithRemovals(t *testing.T) {
	r := UniqueResult{Entries: make([]Entry, 4), Removed: 2}
	summary := FormatUniqueSummary(r)
	if !strings.Contains(summary, "2") || !strings.Contains(summary, "entries") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestFormatUniqueSummary_NoRemovals(t *testing.T) {
	r := UniqueResult{Entries: make([]Entry, 3), Removed: 0}
	summary := FormatUniqueSummary(r)
	if !strings.Contains(summary, "no duplicates") {
		t.Errorf("unexpected summary: %s", summary)
	}
}
