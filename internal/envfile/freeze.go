package envfile

import "fmt"

// FreezeResult holds the outcome of a freeze operation.
type FreezeResult struct {
	Frozen  []string
	Skipped []string
}

// FreezeOptions controls how Freeze behaves.
type FreezeOptions struct {
	// Keys is an explicit list of keys to freeze. If empty, all entries are frozen.
	Keys []string
	// Prefix restricts freezing to entries whose key starts with the given prefix.
	Prefix string
	// DryRun reports what would change without modifying entries.
	DryRun bool
}

// frozenComment is the marker appended to a frozen entry's comment.
const frozenComment = "frozen"

// Freeze marks entries as read-only by appending a "frozen" comment marker.
// Already-frozen entries are skipped.
func Freeze(entries []Entry, opts FreezeOptions) ([]Entry, FreezeResult) {
	keySet := buildKeySet(opts.Keys)
	result := FreezeResult{}
	out := make([]Entry, len(entries))
	copy(out, entries)

	for i, e := range out {
		if e.Key == "" {
			continue
		}
		if opts.Prefix != "" && len(e.Key) < len(opts.Prefix) {
			continue
		}
		if opts.Prefix != "" && e.Key[:len(opts.Prefix)] != opts.Prefix {
			continue
		}
		if len(keySet) > 0 {
			if _, ok := keySet[e.Key]; !ok {
				continue
			}
		}
		if isFrozenEntry(e) {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}
		result.Frozen = append(result.Frozen, e.Key)
		if !opts.DryRun {
			if out[i].Comment == "" {
				out[i].Comment = frozenComment
			} else {
				out[i].Comment = out[i].Comment + " " + frozenComment
			}
		}
	}
	return out, result
}

// isFrozenEntry reports whether an entry carries the frozen marker.
func isFrozenEntry(e Entry) bool {
	for i := 0; i+len(frozenComment) <= len(e.Comment); i++ {
		if e.Comment[i:i+len(frozenComment)] == frozenComment {
			return true
		}
	}
	return false
}

// IsFrozen reports whether the entry with the given key is frozen.
func IsFrozen(entries []Entry, key string) bool {
	for _, e := range entries {
		if e.Key == key {
			return isFrozenEntry(e)
		}
	}
	return false
}

// FormatFreezeSummary returns a human-readable summary of a FreezeResult.
func FormatFreezeSummary(r FreezeResult, dryRun bool) string {
	action := "frozen"
	if dryRun {
		action = "would freeze"
	}
	return fmt.Sprintf("%s %d key(s), %d already frozen",
		action, len(r.Frozen), len(r.Skipped))
}
