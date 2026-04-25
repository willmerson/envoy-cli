package envfile

import (
	"fmt"
	"strings"
)

// FlattenResult holds the outcome of a flatten operation.
type FlattenResult struct {
	Entries  []Entry
	Renamed  int
	Dropped  int
}

// FlattenOptions controls how flattening is performed.
type FlattenOptions struct {
	// Separator replaces the prefix delimiter in key names (default: "_").
	Separator string
	// StripPrefix removes the given prefix from all keys after flattening.
	StripPrefix string
	// UppercaseKeys normalises all keys to uppercase.
	UppercaseKeys bool
	// DropEmpty discards entries whose value is an empty string.
	DropEmpty bool
}

// Flatten normalises a slice of entries by collapsing redundant prefix
// segments, optionally stripping a common prefix and uppercasing keys.
//
// Example: APP_DB__HOST -> APP_DB_HOST (double-separator collapsed to one).
func Flatten(entries []Entry, opts FlattenOptions) FlattenResult {
	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	result := FlattenResult{}
	seen := make(map[string]bool)

	for _, e := range entries {
		if opts.DropEmpty && strings.TrimSpace(e.Value) == "" {
			result.Dropped++
			continue
		}

		key := normaliseKey(e.Key, sep)

		if opts.StripPrefix != "" {
			stripped := strings.TrimPrefix(key, opts.StripPrefix)
			if stripped != key {
				key = strings.TrimPrefix(stripped, sep)
				result.Renamed++
			}
		}

		if opts.UppercaseKeys {
			upper := strings.ToUpper(key)
			if upper != key {
				key = upper
				result.Renamed++
			}
		}

		if seen[key] {
			result.Dropped++
			continue
		}
		seen[key] = true

		result.Entries = append(result.Entries, Entry{
			Key:     key,
			Value:   e.Value,
			Comment: e.Comment,
		})
	}

	return result
}

// FormatFlattenSummary returns a human-readable summary of a FlattenResult.
func FormatFlattenSummary(r FlattenResult) string {
	return fmt.Sprintf("%d entries kept, %d renamed, %d dropped",
		len(r.Entries), r.Renamed, r.Dropped)
}

// normaliseKey collapses consecutive separator characters into one.
func normaliseKey(key, sep string) string {
	double := sep + sep
	for strings.Contains(key, double) {
		key = strings.ReplaceAll(key, double, sep)
	}
	return strings.Trim(key, sep)
}
