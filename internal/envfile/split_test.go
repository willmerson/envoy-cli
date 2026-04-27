package envfile

import (
	"strings"
	"testing"
)

func splitEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET", Value: "abc123"},
	}
}

func TestSplit_ByPrefix(t *testing.T) {
	result := Split(splitEntries(), SplitOptions{ByPrefix: true})

	if len(result.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(result.Groups["DB"]))
	}
	if len(result.Groups["APP"]) != 2 {
		t.Errorf("expected 2 APP entries, got %d", len(result.Groups["APP"]))
	}
	if len(result.Ungrouped) != 1 {
		t.Errorf("expected 1 ungrouped entry, got %d", len(result.Ungrouped))
	}
	if result.Ungrouped[0].Key != "SECRET" {
		t.Errorf("expected ungrouped key SECRET, got %s", result.Ungrouped[0].Key)
	}
}

func TestSplit_LowercaseLabels(t *testing.T) {
	result := Split(splitEntries(), SplitOptions{ByPrefix: true, Lowercase: true})

	if _, ok := result.Groups["db"]; !ok {
		t.Error("expected lowercase group 'db'")
	}
	if _, ok := result.Groups["app"]; !ok {
		t.Error("expected lowercase group 'app'")
	}
}

func TestSplit_CustomSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "DB.HOST", Value: "localhost"},
		{Key: "DB.PORT", Value: "5432"},
		{Key: "STANDALONE", Value: "yes"},
	}
	result := Split(entries, SplitOptions{Separator: "."})

	if len(result.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(result.Groups["DB"]))
	}
	if len(result.Ungrouped) != 1 {
		t.Errorf("expected 1 ungrouped entry, got %d", len(result.Ungrouped))
	}
}

func TestSplit_AllUngrouped(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "1"},
		{Key: "BAR", Value: "2"},
	}
	result := Split(entries, SplitOptions{})

	if len(result.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(result.Groups))
	}
	if len(result.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped, got %d", len(result.Ungrouped))
	}
}

func TestFormatSplitSummary(t *testing.T) {
	result := Split(splitEntries(), SplitOptions{})
	summary := FormatSplitSummary(result)

	if !strings.Contains(summary, "group") {
		t.Error("expected summary to mention groups")
	}
	if !strings.Contains(summary, "ungrouped") {
		t.Error("expected summary to mention ungrouped")
	}
}
