package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// GroupResult holds the grouped entries and metadata.
type GroupResult struct {
	Groups  map[string][]Entry
	Ordered []string // group names in sorted order
}

// GroupOptions controls how grouping is performed.
type GroupOptions struct {
	Separator string // e.g. "_" (default)
	Depth     int    // how many prefix segments to use (default 1)
	Ungrouped string // label for entries with no matching prefix (default "other")
}

// Group partitions entries by key prefix using the given separator and depth.
func Group(entries []Entry, opts GroupOptions) GroupResult {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	if opts.Depth < 1 {
		opts.Depth = 1
	}
	if opts.Ungrouped == "" {
		opts.Ungrouped = "other"
	}

	groups := make(map[string][]Entry)

	for _, e := range entries {
		label := prefixLabel(e.Key, opts.Separator, opts.Depth)
		if label == "" {
			label = opts.Ungrouped
		}
		groups[label] = append(groups[label], e)
	}

	ordered := make([]string, 0, len(groups))
	for k := range groups {
		ordered = append(ordered, k)
	}
	sort.Strings(ordered)

	return GroupResult{Groups: groups, Ordered: ordered}
}

// FormatGroupSummary returns a human-readable summary of the group result.
func FormatGroupSummary(r GroupResult) string {
	var sb strings.Builder
	for _, name := range r.Ordered {
		entries := r.Groups[name]
		sb.WriteString(fmt.Sprintf("[%s] (%d keys)\n", name, len(entries)))
		for _, e := range entries {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Value))
		}
	}
	return sb.String()
}

func prefixLabel(key, sep string, depth int) string {
	parts := strings.Split(key, sep)
	if len(parts) <= depth {
		return ""
	}
	return strings.Join(parts[:depth], sep)
}
