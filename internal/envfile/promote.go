package envfile

import (
	"fmt"
	"strings"
)

// PromoteResult describes what happened during a promote operation.
type PromoteResult struct {
	Added   int
	Updated int
	Skipped int
}

// PromoteOptions controls how promotion behaves.
type PromoteOptions struct {
	// Overwrite existing keys in the target profile
	Overwrite bool
	// PrefixFilter limits promotion to keys with this prefix (empty = all)
	PrefixFilter string
	// StripPrefix removes the prefix from keys when writing to target
	StripPrefix bool
}

// Promote copies entries from src into dst according to opts.
// It returns the merged slice and a summary of changes.
func Promote(src, dst []Entry, opts PromoteOptions) ([]Entry, PromoteResult, error) {
	if src == nil {
		return dst, PromoteResult{}, nil
	}

	dstMap := make(map[string]int, len(dst))
	for i, e := range dst {
		dstMap[e.Key] = i
	}

	result := make([]Entry, len(dst))
	copy(result, dst)

	var res PromoteResult
	for _, e := range src {
		if opts.PrefixFilter != "" && !strings.HasPrefix(e.Key, opts.PrefixFilter) {
			res.Skipped++
			continue
		}

		targetKey := e.Key
		if opts.StripPrefix && opts.PrefixFilter != "" {
			targetKey = strings.TrimPrefix(e.Key, opts.PrefixFilter)
			if targetKey == "" {
				return nil, PromoteResult{}, fmt.Errorf("stripping prefix %q from key %q yields empty key", opts.PrefixFilter, e.Key)
			}
		}

		promoted := Entry{Key: targetKey, Value: e.Value, Comment: e.Comment}

		if idx, exists := dstMap[targetKey]; exists {
			if opts.Overwrite {
				result[idx] = promoted
				res.Updated++
			} else {
				res.Skipped++
			}
		} else {
			dstMap[targetKey] = len(result)
			result = append(result, promoted)
			res.Added++
		}
	}

	return result, res, nil
}

// FormatPromoteResult returns a human-readable summary of a PromoteResult.
func FormatPromoteResult(r PromoteResult) string {
	return fmt.Sprintf("promoted: %d added, %d updated, %d skipped", r.Added, r.Updated, r.Skipped)
}
