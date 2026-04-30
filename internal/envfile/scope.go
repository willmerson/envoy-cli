package envfile

import "strings"

// ScopeOptions controls how Scope filters and transforms entries.
type ScopeOptions struct {
	// Prefix restricts entries to those whose keys begin with this prefix.
	Prefix string
	// StripPrefix removes the prefix from keys in the output.
	StripPrefix bool
	// CaseSensitive controls whether prefix matching is case-sensitive.
	CaseSensitive bool
}

// ScopeResult holds the output of a Scope operation.
type ScopeResult struct {
	Entries []Entry
	Total   int
	Matched int
}

// Scope filters entries to those matching the given prefix and optionally
// strips the prefix from the returned keys.
func Scope(entries []Entry, opts ScopeOptions) ScopeResult {
	result := ScopeResult{Total: len(entries)}

	for _, e := range entries {
		key := e.Key
		prefix := opts.Prefix

		if !opts.CaseSensitive {
			key = strings.ToLower(key)
			prefix = strings.ToLower(prefix)
		}

		if prefix == "" || strings.HasPrefix(key, prefix) {
			out := e
			if opts.StripPrefix && prefix != "" {
				out.Key = e.Key[len(opts.Prefix):]
			}
			result.Entries = append(result.Entries, out)
			result.Matched++
		}
	}

	if result.Entries == nil {
		result.Entries = []Entry{}
	}
	return result
}

// FormatScopeSummary returns a human-readable summary of a ScopeResult.
func FormatScopeSummary(r ScopeResult, prefix string) string {
	if prefix == "" {
		return "scoped all " + itoa(r.Matched) + " entries"
	}
	return "scoped " + itoa(r.Matched) + "/" + itoa(r.Total) +
		" entries matching prefix \"" + prefix + "\""
}
