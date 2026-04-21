package envfile

// DedupeStrategy defines how duplicate keys are resolved.
type DedupeStrategy string

const (
	// DedupeKeepFirst keeps the first occurrence of a duplicate key.
	DedupeKeepFirst DedupeStrategy = "first"
	// DedupeKeepLast keeps the last occurrence of a duplicate key.
	DedupeKeepLast DedupeStrategy = "last"
)

// DedupeSummary holds the result of a deduplication operation.
type DedupeSummary struct {
	Removed int
	Keys    []string
}

// Dedupe removes duplicate keys from a slice of Entry values.
// The strategy determines whether the first or last occurrence is kept.
func Dedupe(entries []Entry, strategy DedupeStrategy) ([]Entry, DedupeSummary) {
	seen := make(map[string]int) // key -> index in result
	result := make([]Entry, 0, len(entries))
	summary := DedupeSummary{}

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			switch strategy {
			case DedupeKeepLast:
				// Replace the existing entry with the new one.
				result[idx] = e
				summary.Removed++
				summary.Keys = append(summary.Keys, e.Key)
			default: // DedupeKeepFirst
				// Discard the new entry.
				summary.Removed++
				summary.Keys = append(summary.Keys, e.Key)
			}
		} else {
			seen[e.Key] = len(result)
			result = append(result, e)
		}
	}

	return result, summary
}
