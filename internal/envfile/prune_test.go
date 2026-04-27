package envfile

import (
	"strings"
	"testing"
)

func pruneEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_DEBUG", Value: ""},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: ""},
		{Key: "SECRET_KEY", Value: "abc123"},
		{Comment: true, Raw: "# legacy setting"},
	}
}

func TestPrune_RemoveEmpty(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{RemoveEmpty: true})
	for _, e := range result.Entries {
		if !e.Comment && e.Value == "" {
			t.Errorf("expected empty entry %q to be removed", e.Key)
		}
	}
	if len(result.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestPrune_RemoveCommented(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{RemoveCommented: true})
	for _, e := range result.Entries {
		if e.Comment {
			t.Error("expected comment entries to be removed")
		}
	}
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed))
	}
}

func TestPrune_RemoveKeys(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{RemoveKeys: []string{"APP_NAME", "SECRET_KEY"}})
	for _, e := range result.Entries {
		if e.Key == "APP_NAME" || e.Key == "SECRET_KEY" {
			t.Errorf("expected key %q to be removed", e.Key)
		}
	}
	if len(result.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestPrune_RemovePrefix(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{RemovePrefix: "DB_"})
	for _, e := range result.Entries {
		if strings.HasPrefix(e.Key, "DB_") {
			t.Errorf("expected key %q to be pruned by prefix", e.Key)
		}
	}
	if len(result.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestPrune_NothingRemoved(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{})
	if len(result.Removed) != 0 {
		t.Errorf("expected nothing removed, got %d", len(result.Removed))
	}
	if len(result.Entries) != len(pruneEntries()) {
		t.Errorf("expected all entries kept")
	}
}

func TestFormatPruneSummary_NothingRemoved(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{})
	summary := FormatPruneSummary(result)
	if !strings.Contains(summary, "nothing removed") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestFormatPruneSummary_WithRemovals(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{RemovePrefix: "DB_"})
	summary := FormatPruneSummary(result)
	if !strings.Contains(summary, "DB_HOST") {
		t.Errorf("expected DB_HOST in summary, got: %s", summary)
	}
}
