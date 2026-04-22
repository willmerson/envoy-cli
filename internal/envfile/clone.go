package envfile

import (
	"fmt"
	"slices"
)

// CloneResult holds the result of a clone operation.
type CloneResult struct {
	Copied  int
	Skipped int
	Keys    []string
}

// CloneOptions controls how Clone behaves.
type CloneOptions struct {
	// Prefix filters source entries to only those whose key starts with Prefix.
	// If empty, all entries are cloned.
	Prefix string
	// DestPrefix replaces the source prefix in cloned keys.
	// Only applied when Prefix is non-empty.
	DestPrefix string
	// Overwrite controls whether existing keys in dst are overwritten.
	Overwrite bool
}

// Clone copies entries from src into dst according to opts.
// It returns a CloneResult describing what was copied or skipped.
func Clone(src, dst []Entry, opts CloneOptions) ([]Entry, CloneResult, error) {
	result := CloneResult{}
	out := make([]Entry, len(dst))
	copy(out, dst)

	// Build a lookup set of existing keys in dst.
	existing := make(map[string]int, len(dst))
	for i, e := range out {
		existing[e.Key] = i
	}

	for _, e := range src {
		key := e.Key

		// Apply prefix filter.
		if opts.Prefix != "" {
			if len(key) < len(opts.Prefix) || key[:len(opts.Prefix)] != opts.Prefix {
				continue
			}
			if opts.DestPrefix != "" {
				key = opts.DestPrefix + key[len(opts.Prefix):]
			}
		}

		if key == "" {
			return nil, result, fmt.Errorf("clone produced empty key from source key %q", e.Key)
		}

		newEntry := Entry{Key: key, Value: e.Value, Comment: e.Comment}

		if idx, exists := existing[key]; exists {
			if !opts.Overwrite {
				result.Skipped++
				continue
			}
			out[idx] = newEntry
		} else {
			out = append(out, newEntry)
			existing[key] = len(out) - 1
		}

		result.Copied++
		if !slices.Contains(result.Keys, key) {
			result.Keys = append(result.Keys, key)
		}
	}

	return out, result, nil
}
