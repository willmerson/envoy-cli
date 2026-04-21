package envfile

import (
	"strings"
)

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	Prefix    string
	Suffix    string
	Contains  string
	KeysOnly  bool
	EmptyOnly bool
}

// FilterResult holds the result of a filter operation.
type FilterResult struct {
	Matched []Entry
	Total   int
}

// Filter returns entries matching the given options.
func Filter(entries []Entry, opts FilterOptions) FilterResult {
	result := FilterResult{
		Total:   len(entries),
		Matched: []Entry{},
	}

	for _, e := range entries {
		if !matchesFilter(e, opts) {
			continue
		}
		result.Matched = append(result.Matched, e)
	}

	return result
}

// Count returns the number of entries that match the given options.
func Count(entries []Entry, opts FilterOptions) int {
	return Filter(entries, opts).MatchedCount()
}

// MatchedCount returns the number of matched entries.
func (r FilterResult) MatchedCount() int {
	return len(r.Matched)
}

func matchesFilter(e Entry, opts FilterOptions) bool {
	if opts.EmptyOnly && e.Value != "" {
		return false
	}

	target := e.Key
	if !opts.KeysOnly {
		target = e.Key + "=" + e.Value
	}

	if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
		return false
	}

	if opts.Suffix != "" && !strings.HasSuffix(e.Key, opts.Suffix) {
		return false
	}

	if opts.Contains != "" && !strings.Contains(target, opts.Contains) {
		return false
	}

	return true
}
