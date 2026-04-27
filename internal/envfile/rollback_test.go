package envfile

import (
	"testing"
)

func rbEntries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestRollback_RestoresChangedValues(t *testing.T) {
	current := rbEntries("DB_HOST", "new-host", "APP_ENV", "production")
	snap := rbEntries("DB_HOST", "old-host", "APP_ENV", "production")

	restored, result := Rollback(current, snap)
	if result.Restored != 1 {
		t.Errorf("expected 1 restored, got %d", result.Restored)
	}
	if restored[0].Value != "old-host" {
		t.Errorf("expected old-host, got %s", restored[0].Value)
	}
}

func TestRollback_AddsKeysFromSnapshot(t *testing.T) {
	current := rbEntries("APP_ENV", "staging")
	snap := rbEntries("APP_ENV", "staging", "DB_HOST", "localhost")

	_, result := Rollback(current, snap)
	if result.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Added)
	}
}

func TestRollback_RemovesExtraKeys(t *testing.T) {
	current := rbEntries("APP_ENV", "staging", "EXTRA", "value")
	snap := rbEntries("APP_ENV", "staging")

	_, result := Rollback(current, snap)
	if result.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Removed)
	}
}

func TestRollback_NoChanges(t *testing.T) {
	entries := rbEntries("A", "1", "B", "2")
	_, result := Rollback(entries, entries)
	if result.Restored != 0 || result.Added != 0 || result.Removed != 0 {
		t.Errorf("expected no changes, got %+v", result)
	}
}

func TestFormatRollbackResult(t *testing.T) {
	r := RollbackResult{SnapshotName: "snap-1", Restored: 2, Added: 1, Removed: 3}
	msg := FormatRollbackResult(r)
	if msg == "" {
		t.Error("expected non-empty format string")
	}
}

func TestPlanRollback(t *testing.T) {
	current := rbEntries("A", "old", "B", "same", "C", "extra")
	snap := rbEntries("A", "new", "B", "same", "D", "added")

	plan := PlanRollback(current, snap)
	if len(plan.ToRestore) != 1 || plan.ToRestore[0] != "A" {
		t.Errorf("unexpected ToRestore: %v", plan.ToRestore)
	}
	if len(plan.ToAdd) != 1 || plan.ToAdd[0] != "D" {
		t.Errorf("unexpected ToAdd: %v", plan.ToAdd)
	}
	if len(plan.ToRemove) != 1 || plan.ToRemove[0] != "C" {
		t.Errorf("unexpected ToRemove: %v", plan.ToRemove)
	}
}
