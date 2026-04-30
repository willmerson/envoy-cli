package envfile

import (
	"strings"
	"testing"
)

func freezeEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost", Comment: "frozen"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestFreeze_ExplicitKeys(t *testing.T) {
	entries := freezeEntries()
	out, result := Freeze(entries, FreezeOptions{Keys: []string{"APP_NAME", "APP_ENV"}})
	if len(result.Frozen) != 2 {
		t.Fatalf("expected 2 frozen, got %d", len(result.Frozen))
	}
	if !isFrozenEntry(out[0]) || !isFrozenEntry(out[1]) {
		t.Error("expected APP_NAME and APP_ENV to be frozen")
	}
}

func TestFreeze_SkipsAlreadyFrozen(t *testing.T) {
	entries := freezeEntries()
	_, result := Freeze(entries, FreezeOptions{Keys: []string{"DB_HOST"}})
	if len(result.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if len(result.Frozen) != 0 {
		t.Fatalf("expected 0 frozen, got %d", len(result.Frozen))
	}
}

func TestFreeze_PrefixFilter(t *testing.T) {
	entries := freezeEntries()
	_, result := Freeze(entries, FreezeOptions{Prefix: "APP_"})
	if len(result.Frozen) != 2 {
		t.Fatalf("expected 2 frozen via prefix, got %d", len(result.Frozen))
	}
}

func TestFreeze_AllKeys(t *testing.T) {
	entries := freezeEntries()
	out, result := Freeze(entries, FreezeOptions{})
	// DB_HOST already frozen => skipped
	if len(result.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if len(result.Frozen) != 4 {
		t.Fatalf("expected 4 frozen, got %d", len(result.Frozen))
	}
	for _, e := range out {
		if e.Key == "" {
			continue
		}
		if !isFrozenEntry(e) {
			t.Errorf("expected %s to be frozen", e.Key)
		}
	}
}

func TestFreeze_DryRun(t *testing.T) {
	entries := freezeEntries()
	out, result := Freeze(entries, FreezeOptions{Keys: []string{"APP_NAME"}, DryRun: true})
	if len(result.Frozen) != 1 {
		t.Fatalf("expected 1 in dry-run frozen list, got %d", len(result.Frozen))
	}
	// Original entry must not be modified
	if isFrozenEntry(out[0]) {
		t.Error("dry-run should not modify entries")
	}
}

func TestIsFrozen(t *testing.T) {
	entries := freezeEntries()
	if !IsFrozen(entries, "DB_HOST") {
		t.Error("DB_HOST should be frozen")
	}
	if IsFrozen(entries, "APP_NAME") {
		t.Error("APP_NAME should not be frozen")
	}
}

func TestFormatFreezeSummary(t *testing.T) {
	r := FreezeResult{Frozen: []string{"A", "B"}, Skipped: []string{"C"}}
	s := FormatFreezeSummary(r, false)
	if !strings.Contains(s, "2") || !strings.Contains(s, "1") {
		t.Errorf("unexpected summary: %s", s)
	}
	ds := FormatFreezeSummary(r, true)
	if !strings.Contains(ds, "would freeze") {
		t.Errorf("expected dry-run wording, got: %s", ds)
	}
}
