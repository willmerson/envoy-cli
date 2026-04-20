package envfile

import (
	"fmt"
	"strings"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Found   bool
	Updated bool
}

// RenameKey renames a key in a slice of Entry values.
// Returns the updated entries and a RenameResult describing what happened.
func RenameKey(entries []Entry, oldKey, newKey string) ([]Entry, RenameResult, error) {
	if strings.TrimSpace(oldKey) == "" {
		return entries, RenameResult{}, fmt.Errorf("old key must not be empty")
	}
	if strings.TrimSpace(newKey) == "" {
		return entries, RenameResult{}, fmt.Errorf("new key must not be empty")
	}
	if oldKey == newKey {
		return entries, RenameResult{OldKey: oldKey, NewKey: newKey, Found: false, Updated: false},
			fmt.Errorf("old key and new key are identical")
	}

	// Check if newKey already exists
	for _, e := range entries {
		if e.Key == newKey {
			return entries, RenameResult{}, fmt.Errorf("key %q already exists", newKey)
		}
	}

	result := RenameResult{OldKey: oldKey, NewKey: newKey}
	updated := make([]Entry, len(entries))
	copy(updated, entries)

	for i, e := range updated {
		if e.Key == oldKey {
			updated[i].Key = newKey
			result.Found = true
			result.Updated = true
			break
		}
	}

	if !result.Found {
		return entries, result, fmt.Errorf("key %q not found", oldKey)
	}

	return updated, result, nil
}
