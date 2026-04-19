package envfile

import "fmt"

// DiffStatus represents the type of change between two env files.
type DiffStatus string

const (
	DiffAdded   DiffStatus = "added"
	DiffRemoved DiffStatus = "removed"
	DiffChanged DiffStatus = "changed"
	DiffUnchanged DiffStatus = "unchanged"
)

// DiffEntry represents a single key's diff result.
type DiffEntry struct {
	Key      string
	OldValue string
	NewValue string
	Status   DiffStatus
}

// Diff compares two slices of Entry and returns a list of DiffEntry.
func Diff(base, other []Entry) []DiffEntry {
	baseMap := ToMap(base)
	otherMap := ToMap(other)

	seen := map[string]bool{}
	var results []DiffEntry

	for _, e := range base {
		seen[e.Key] = true
		newVal, exists := otherMap[e.Key]
		if !exists {
			results = append(results, DiffEntry{Key: e.Key, OldValue: e.Value, Status: DiffRemoved})
		} else if newVal != e.Value {
			results = append(results, DiffEntry{Key: e.Key, OldValue: e.Value, NewValue: newVal, Status: DiffChanged})
		} else {
			results = append(results, DiffEntry{Key: e.Key, OldValue: e.Value, NewValue: newVal, Status: DiffUnchanged})
		}
	}

	for _, e := range other {
		if !seen[e.Key] {
			results = append(results, DiffEntry{Key: e.Key, NewValue: e.Value, Status: DiffAdded})
		}
	}

	return results
}

// FormatDiff returns a human-readable string of the diff.
func FormatDiff(entries []DiffEntry) string {
	var out string
	for _, d := range entries {
		switch d.Status {
		case DiffAdded:
			out += fmt.Sprintf("+ %s=%s\n", d.Key, d.NewValue)
		case DiffRemoved:
			out += fmt.Sprintf("- %s=%s\n", d.Key, d.OldValue)
		case DiffChanged:
			out += fmt.Sprintf("~ %s: %s -> %s\n", d.Key, d.OldValue, d.NewValue)
		}
	}
	return out
}
