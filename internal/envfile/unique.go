package envfile

// UniqueResult holds the outcome of a Unique operation.
type UniqueResult struct {
	Entries  []Entry
	Removed  int
	KeptKeys []string
}

// UniqueOptions controls how Unique behaves.
type UniqueOptions struct {
	// ByValue removes entries whose value is duplicated, keeping the first occurrence.
	ByValue bool
	// ByKey removes entries whose key is duplicated (alias for Dedupe KeepFirst).
	ByKey bool
	// CaseSensitive controls whether value comparison is case-sensitive.
	CaseSensitive bool
}

// Unique filters entries so that each key (and optionally each value) appears
// only once. It returns a UniqueResult describing what was kept and removed.
func Unique(entries []Entry, opts UniqueOptions) UniqueResult {
	seenKeys := make(map[string]bool)
	seenValues := make(map[string]bool)

	var kept []Entry
	removed := 0

	for _, e := range entries {
		keyNorm := e.Key
		valNorm := e.Value
		if !opts.CaseSensitive {
			valNorm = toLowerStr(e.Value)
		}

		if opts.ByKey && seenKeys[keyNorm] {
			removed++
			continue
		}
		if opts.ByValue && seenValues[valNorm] {
			removed++
			continue
		}

		seenKeys[keyNorm] = true
		seenValues[valNorm] = true
		kept = append(kept, e)
	}

	var keptKeys []string
	for _, e := range kept {
		keptKeys = append(keptKeys, e.Key)
	}

	return UniqueResult{
		Entries:  kept,
		Removed:  removed,
		KeptKeys: keptKeys,
	}
}

// FormatUniqueSummary returns a human-readable summary of a UniqueResult.
func FormatUniqueSummary(r UniqueResult) string {
	if r.Removed == 0 {
		return "unique: no duplicates found"
	}
	suffix := "entry"
	if r.Removed != 1 {
		suffix = "entries"
	}
	return fmt.Sprintf("unique: removed %d duplicate %s, %d remaining", r.Removed, suffix, len(r.Entries))
}

func toLowerStr(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
