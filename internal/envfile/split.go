package envfile

import "fmt"

// SplitOptions controls how entries are split into multiple files.
type SplitOptions struct {
	// ByPrefix splits entries into groups based on key prefix (e.g. "DB_" -> db.env)
	ByPrefix bool
	// Separator is the delimiter used to detect prefixes (default "_")
	Separator string
	// Lowercase controls whether output filenames are lowercased
	Lowercase bool
}

// SplitResult holds the output of a Split operation.
type SplitResult struct {
	// Groups maps a label (e.g. "DB") to its entries
	Groups map[string][]Entry
	// Ungrouped holds entries that did not match any prefix group
	Ungrouped []Entry
}

// Split partitions entries into named groups based on key prefix.
// Entries whose keys have no prefix separator are placed in Ungrouped.
func Split(entries []Entry, opts SplitOptions) SplitResult {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	groups := make(map[string][]Entry)
	var ungrouped []Entry

	for _, e := range entries {
		label, ok := prefixOf(e.Key, opts.Separator)
		if !ok {
			ungrouped = append(ungrouped, e)
			continue
		}
		if opts.Lowercase {
			label = toLower(label)
		}
		groups[label] = append(groups[label], e)
	}

	return SplitResult{Groups: groups, Ungrouped: ungrouped}
}

// FormatSplitSummary returns a human-readable summary of a SplitResult.
func FormatSplitSummary(r SplitResult) string {
	out := fmt.Sprintf("Split into %d group(s):\n", len(r.Groups))
	for label, entries := range r.Groups {
		out += fmt.Sprintf("  %-20s %d key(s)\n", label, len(entries))
	}
	if len(r.Ungrouped) > 0 {
		out += fmt.Sprintf("  %-20s %d key(s)\n", "(ungrouped)", len(r.Ungrouped))
	}
	return out
}

func prefixOf(key, sep string) (string, bool) {
	for i := 0; i < len(key)-1; i++ {
		if string(key[i]) == sep {
			return key[:i], true
		}
	}
	return "", false
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
