package envfile

import (
	"testing"
)

func entries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMerge_AddsNewKeys(t *testing.T) {
	base := entries("A", "1")
	override := entries("B", "2")
	r := Merge(base, override, StrategyOurs)
	if len(r.Added) != 1 || r.Added[0] != "B" {
		t.Errorf("expected B added, got %v", r.Added)
	}
	if len(r.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestMerge_ConflictStrategyOurs(t *testing.T) {
	base := entries("A", "original")
	override := entries("A", "new")
	r := Merge(base, override, StrategyOurs)
	if r.Entries[0].Value != "original" {
		t.Errorf("expected original value kept, got %s", r.Entries[0].Value)
	}
	if len(r.Overridden) != 0 {
		t.Errorf("expected no overrides with StrategyOurs")
	}
}

func TestMerge_ConflictStrategyTheirs(t *testing.T) {
	base := entries("A", "original")
	override := entries("A", "new")
	r := Merge(base, override, StrategyTheirs)
	if r.Entries[0].Value != "new" {
		t.Errorf("expected new value, got %s", r.Entries[0].Value)
	}
	if len(r.Overridden) != 1 {
		t.Errorf("expected 1 override, got %d", len(r.Overridden))
	}
}

func TestMergeSummary(t *testing.T) {
	r := MergeResult{Added: []string{"X"}, Conflicts: []string{"Y"}, Overridden: []string{"Y"}}
	s := MergeSummary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
