package envfile

import "fmt"

// SquashResult holds the outcome of a squash operation.
type SquashResult struct {
	Entries  []Entry
	Kept     int
	Squashed int
}

// SquashOptions controls how Squash behaves.
type SquashOptions struct {
	// PrefixGroups, when true, squashes only within each prefix group.
	// When false, squashes across all entries by key.
	PrefixGroups bool
	// Separator is the prefix delimiter (default "_").
	Separator string
	// KeepLast, when true, keeps the last occurrence; otherwise keeps the first.
	KeepLast bool
}

// Squash deduplicates entries by key, optionally scoped to prefix groups.
// It differs from Dedupe in that it can operate within prefix groups and
// always returns a SquashResult with change counts.
func Squash(entries []Entry, opts SquashOptions) SquashResult {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	if !opts.PrefixGroups {
		return squashFlat(entries, opts.KeepLast)
	}
	return squashByPrefix(entries, opts.Separator, opts.KeepLast)
}

func squashFlat(entries []Entry, keepLast bool) SquashResult {
	seen := make(map[string]int) // key -> index in out
	out := make([]Entry, 0, len(entries))
	squashed := 0

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			if keepLast {
				out[idx] = e
			}
			squashed++
		} else {
			seen[e.Key] = len(out)
			out = append(out, e)
		}
	}

	return SquashResult{Entries: out, Kept: len(out), Squashed: squashed}
}

func squashByPrefix(entries []Entry, sep string, keepLast bool) SquashResult {
	// Group entries by prefix, preserving insertion order of groups.
	type group struct {
		label   string
		indices []int
	}
	groupOrder := []string{}
	groupMap := map[string]*group{}

	for i, e := range entries {
		label := prefixLabel(e.Key, sep)
		if _, ok := groupMap[label]; !ok {
			groupMap[label] = &group{label: label}
			groupOrder = append(groupOrder, label)
		}
		groupMap[label].indices = append(groupMap[label].indices, i)
	}

	out := make([]Entry, 0, len(entries))
	squashed := 0

	for _, label := range groupOrder {
		g := groupMap[label]
		groupEntries := make([]Entry, 0, len(g.indices))
		for _, idx := range g.indices {
			groupEntries = append(groupEntries, entries[idx])
		}
		res := squashFlat(groupEntries, keepLast)
		out = append(out, res.Entries...)
		squashed += res.Squashed
	}

	return SquashResult{Entries: out, Kept: len(out), Squashed: squashed}
}

// FormatSquashSummary returns a human-readable summary of a SquashResult.
func FormatSquashSummary(r SquashResult) string {
	if r.Squashed == 0 {
		return "squash: no duplicate keys found"
	}
	return fmt.Sprintf("squash: removed %d duplicate key(s), %d unique key(s) kept", r.Squashed, r.Kept)
}
