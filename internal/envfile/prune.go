package envfile

import "strings"

// PruneOptions controls which entries are removed.
type PruneOptions struct {
	// RemoveEmpty removes entries whose value is an empty string.
	RemoveEmpty bool
	// RemoveCommented removes entries that are commented out (raw comment lines).
	RemoveCommented bool
	// RemoveKeys removes entries whose key exactly matches one of these values.
	RemoveKeys []string
	// RemovePrefix removes entries whose key starts with the given prefix.
	RemovePrefix string
}

// PruneResult holds the outcome of a Prune operation.
type PruneResult struct {
	Entries []Entry
	Removed []Entry
}

// FormatPruneSummary returns a human-readable summary of the prune result.
func FormatPruneSummary(r PruneResult) string {
	if len(r.Removed) == 0 {
		return "prune: nothing removed"
	}
	var sb strings.Builder
	sb.WriteString("prune: removed entries\n")
	for _, e := range r.Removed {
		if e.Comment {
			sb.WriteString("  # (comment line)\n")
		} else {
			sb.WriteString("  " + e.Key + "\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Prune removes entries from the list according to PruneOptions and returns
// the kept entries together with the removed ones.
func Prune(entries []Entry, opts PruneOptions) PruneResult {
	removeKeySet := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		removeKeySet[k] = struct{}{}
	}

	var kept, removed []Entry
	for _, e := range entries {
		switch {
		case opts.RemoveCommented && e.Comment:
			removed = append(removed, e)
		case opts.RemoveEmpty && !e.Comment && e.Value == "":
			removed = append(removed, e)
		case !e.Comment && len(removeKeySet) > 0:
			if _, ok := removeKeySet[e.Key]; ok {
				removed = append(removed, e)
			} else {
				kept = append(kept, e)
			}
		case !e.Comment && opts.RemovePrefix != "" && strings.HasPrefix(e.Key, opts.RemovePrefix):
			removed = append(removed, e)
		default:
			kept = append(kept, e)
		}
	}
	return PruneResult{Entries: kept, Removed: removed}
}
