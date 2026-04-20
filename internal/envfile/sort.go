package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the order in which entries should be sorted.
type SortOrder int

const (
	SortAsc  SortOrder = iota // alphabetical ascending
	SortDesc                  // alphabetical descending
)

// SortResult holds the outcome of a sort operation.
type SortResult struct {
	Entries []Entry
	Moved   int // number of entries that changed position
}

// Sort returns a new slice of entries sorted by key according to the given order.
// Comments and blank entries are moved to the top, preserving their relative order.
func Sort(entries []Entry, order SortOrder) SortResult {
	var comments []Entry
	var keyed []Entry

	for _, e := range entries {
		if strings.HasPrefix(strings.TrimSpace(e.Key), "#") || e.Key == "" {
			comments = append(comments, e)
		} else {
			keyed = append(keyed, e)
		}
	}

	sorted := make([]Entry, len(keyed))
	copy(sorted, keyed)

	sort.SliceStable(sorted, func(i, j int) bool {
		ki := strings.ToLower(sorted[i].Key)
		kj := strings.ToLower(sorted[j].Key)
		if order == SortDesc {
			return ki > kj
		}
		return ki < kj
	})

	// count moved entries
	moved := 0
	for i, e := range sorted {
		if i < len(keyed) && keyed[i].Key != e.Key {
			moved++
		}
	}

	result := append(comments, sorted...)
	return SortResult{Entries: result, Moved: moved}
}

// GroupByPrefix groups entries by their key prefix (split on "_").
// Entries without an underscore are placed under the "" (empty) group.
func GroupByPrefix(entries []Entry) map[string][]Entry {
	groups := make(map[string][]Entry)
	for _, e := range entries {
		parts := strings.SplitN(e.Key, "_", 2)
		prefix := ""
		if len(parts) == 2 {
			prefix = parts[0]
		}
		groups[prefix] = append(groups[prefix], e)
	}
	return groups
}
