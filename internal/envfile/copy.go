package envfile

import (
	"fmt"
)

// CopyResult holds the result of a copy operation.
type CopyResult struct {
	Key      string
	OldValue string
	NewKey   string
}

// CopyKey duplicates the value of an existing key into a new key.
// Returns an error if oldKey does not exist or newKey already exists.
func CopyKey(entries []Entry, oldKey, newKey string) ([]Entry, CopyResult, error) {
	if oldKey == "" {
		return nil, CopyResult{}, fmt.Errorf("old key must not be empty")
	}
	if newKey == "" {
		return nil, CopyResult{}, fmt.Errorf("new key must not be empty")
	}
	if oldKey == newKey {
		return nil, CopyResult{}, fmt.Errorf("old key and new key must be different")
	}

	var foundEntry *Entry
	for i := range entries {
		if entries[i].Key == oldKey {
			foundEntry = &entries[i]
		}
		if entries[i].Key == newKey {
			return nil, CopyResult{}, fmt.Errorf("key %q already exists", newKey)
		}
	}

	if foundEntry == nil {
		return nil, CopyResult{}, fmt.Errorf("key %q not found", oldKey)
	}

	newEntry := Entry{
		Key:     newKey,
		Value:   foundEntry.Value,
		Comment: fmt.Sprintf("copied from %s", oldKey),
	}

	result := make([]Entry, len(entries)+1)
	copy(result, entries)
	result[len(entries)] = newEntry

	return result, CopyResult{
		Key:      oldKey,
		OldValue: foundEntry.Value,
		NewKey:   newKey,
	}, nil
}
