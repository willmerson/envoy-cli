package envfile

import (
	"strings"
	"testing"
)

func auditEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestAudit_Added(t *testing.T) {
	before := auditEntries("FOO", "bar")
	after := auditEntries("FOO", "bar", "NEW_KEY", "hello")
	log := Audit(before, after)
	if len(log) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(log))
	}
	if log[0].Action != AuditAdded || log[0].Key != "NEW_KEY" || log[0].NewValue != "hello" {
		t.Errorf("unexpected audit entry: %+v", log[0])
	}
}

func TestAudit_Removed(t *testing.T) {
	before := auditEntries("FOO", "bar", "OLD_KEY", "bye")
	after := auditEntries("FOO", "bar")
	log := Audit(before, after)
	if len(log) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(log))
	}
	if log[0].Action != AuditRemoved || log[0].Key != "OLD_KEY" || log[0].OldValue != "bye" {
		t.Errorf("unexpected audit entry: %+v", log[0])
	}
}

func TestAudit_Changed(t *testing.T) {
	before := auditEntries("FOO", "old")
	after := auditEntries("FOO", "new")
	log := Audit(before, after)
	if len(log) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(log))
	}
	if log[0].Action != AuditChanged || log[0].OldValue != "old" || log[0].NewValue != "new" {
		t.Errorf("unexpected audit entry: %+v", log[0])
	}
}

func TestAudit_NoChanges(t *testing.T) {
	before := auditEntries("FOO", "bar", "BAZ", "qux")
	after := auditEntries("FOO", "bar", "BAZ", "qux")
	log := Audit(before, after)
	if len(log) != 0 {
		t.Errorf("expected no audit entries, got %d", len(log))
	}
}

func TestFormatAuditLog_Empty(t *testing.T) {
	out := FormatAuditLog(AuditLog{})
	if out != "No changes detected." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatAuditLog_ContainsActions(t *testing.T) {
	before := auditEntries("A", "1", "B", "2")
	after := auditEntries("A", "99", "C", "3")
	log := Audit(before, after)
	out := FormatAuditLog(log)
	if !strings.Contains(out, "ADDED") {
		t.Error("expected ADDED in output")
	}
	if !strings.Contains(out, "REMOVED") {
		t.Error("expected REMOVED in output")
	}
	if !strings.Contains(out, "CHANGED") {
		t.Error("expected CHANGED in output")
	}
}
