package envfile

import "fmt"

// MergeStrategy defines how conflicts are resolved during merge.
type MergeStrategy int

const (
	// StrategyOurs keeps the base value on conflict.
	StrategyOurs MergeStrategy = iota
	// StrategyTheirs overwrites with the incoming value on conflict.
	StrategyTheirs
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Entries    []Entry
	Conflicts  []string
	Added      []string
	Overridden []string
}

// Merge combines base entries with override entries using the given strategy.
// It returns a MergeResult describing what changed.
func Merge(base, override []Entry, strategy MergeStrategy) MergeResult {
	result := MergeResult{}
	index := make(map[string]int, len(base))
	merged := make([]Entry, len(base))
	copy(merged, base)

	for i, e := range merged {
		index[e.Key] = i
	}

	for _, e := range override {
		if e.Comment {
			continue
		}
		if idx, exists := index[e.Key]; exists {
			result.Conflicts = append(result.Conflicts, e.Key)
			if strategy == StrategyTheirs {
				merged[idx].Value = e.Value
				result.Overridden = append(result.Overridden, e.Key)
			}
		} else {
			merged = append(merged, e)
			index[e.Key] = len(merged) - 1
			result.Added = append(result.Added, e.Key)
		}
	}

	result.Entries = merged
	return result
}

// MergeSummary returns a human-readable summary of a MergeResult.
func MergeSummary(r MergeResult) string {
	return fmt.Sprintf("added: %d, conflicts: %d, overridden: %d",
		len(r.Added), len(r.Conflicts), len(r.Overridden))
}
