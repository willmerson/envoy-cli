package envfile

import (
	"fmt"
	"strings"
	"time"
)

// AuditAction represents the type of change made to an env entry.
type AuditAction string

const (
	AuditAdded   AuditAction = "added"
	AuditRemoved AuditAction = "removed"
	AuditChanged AuditAction = "changed"
)

// AuditEntry records a single change event.
type AuditEntry struct {
	Timestamp time.Time
	Action    AuditAction
	Key       string
	OldValue  string
	NewValue  string
}

// AuditLog holds a list of audit entries.
type AuditLog []AuditEntry

// Audit compares two slices of Entry and returns an AuditLog describing
// what was added, removed, or changed between the before and after states.
func Audit(before, after []Entry) AuditLog {
	now := time.Now().UTC()
	beforeMap := make(map[string]string, len(before))
	for _, e := range before {
		beforeMap[e.Key] = e.Value
	}

	afterMap := make(map[string]string, len(after))
	for _, e := range after {
		afterMap[e.Key] = e.Value
	}

	var log AuditLog

	for _, e := range after {
		oldVal, existed := beforeMap[e.Key]
		if !existed {
			log = append(log, AuditEntry{Timestamp: now, Action: AuditAdded, Key: e.Key, NewValue: e.Value})
		} else if oldVal != e.Value {
			log = append(log, AuditEntry{Timestamp: now, Action: AuditChanged, Key: e.Key, OldValue: oldVal, NewValue: e.Value})
		}
	}

	for _, e := range before {
		if _, exists := afterMap[e.Key]; !exists {
			log = append(log, AuditEntry{Timestamp: now, Action: AuditRemoved, Key: e.Key, OldValue: e.Value})
		}
	}

	return log
}

// FormatAuditLog returns a human-readable string representation of the audit log.
func FormatAuditLog(log AuditLog) string {
	if len(log) == 0 {
		return "No changes detected."
	}
	var sb strings.Builder
	for _, entry := range log {
		ts := entry.Timestamp.Format(time.RFC3339)
		switch entry.Action {
		case AuditAdded:
			sb.WriteString(fmt.Sprintf("[%s] ADDED    %s = %q\n", ts, entry.Key, entry.NewValue))
		case AuditRemoved:
			sb.WriteString(fmt.Sprintf("[%s] REMOVED  %s (was %q)\n", ts, entry.Key, entry.OldValue))
		case AuditChanged:
			sb.WriteString(fmt.Sprintf("[%s] CHANGED  %s: %q -> %q\n", ts, entry.Key, entry.OldValue, entry.NewValue))
		}
	}
	return sb.String()
}
