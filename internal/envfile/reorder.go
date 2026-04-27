package envfile

import (
	"fmt"
	"strings"
)

// ReorderOptions controls how entries are reordered.
type ReorderOptions struct {
	// Keys defines the explicit order for known keys.
	// Keys not listed here are placed after, preserving their relative order.
	Keys []string
	// PinFirst moves all entries whose keys start with these prefixes to the top.
	PinFirst []string
}

// ReorderResult holds the outcome of a Reorder operation.
type ReorderResult struct {
	Entries  []Entry
	Moved    int
	Unknown  int
}

// Reorder rearranges entries according to the provided options.
// Explicitly listed keys come first (in specified order), followed by
// prefix-pinned entries, then any remaining entries in original order.
func Reorder(entries []Entry, opts ReorderOptions) ReorderResult {
	indexed := make(map[string]Entry, len(entries))
	for _, e := range entries {
		indexed[e.Key] = e
	}

	seen := make(map[string]bool)
	var result []Entry

	// 1. Explicit key order
	for _, k := range opts.Keys {
		if e, ok := indexed[k]; ok {
			result = append(result, e)
			seen[k] = true
		}
	}

	// 2. Prefix-pinned entries not already placed
	for _, e := range entries {
		if seen[e.Key] {
			continue
		}
		for _, pfx := range opts.PinFirst {
			if strings.HasPrefix(e.Key, pfx) {
				result = append(result, e)
				seen[e.Key] = true
				break
			}
		}
	}

	// 3. Remaining entries in original order
	unknown := 0
	for _, e := range entries {
		if !seen[e.Key] {
			result = append(result, e)
			unknown++
		}
	}

	moved := len(entries) - unknown

	return ReorderResult{
		Entries: result,
		Moved:   moved,
		Unknown: unknown,
	}
}

// FormatReorderSummary returns a human-readable summary of a ReorderResult.
func FormatReorderSummary(r ReorderResult) string {
	return fmt.Sprintf("%d entr%s repositioned, %d left in original order",
		r.Moved, pluralSuffix(r.Moved), r.Unknown)
}

func pluralSuffix(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
