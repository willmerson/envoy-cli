package envfile

import (
	"fmt"
	"strings"
)

// MaskOptions controls how values are masked.
type MaskOptions struct {
	// Keys is a list of exact key names to mask.
	Keys []string
	// Patterns is a list of substring patterns; any key containing a pattern is masked.
	Patterns []string
	// Placeholder replaces the visible portion. Defaults to "****".
	Placeholder string
	// ShowLast reveals the last N characters of the value. 0 means show nothing.
	ShowLast int
}

// MaskResult holds the outcome of a Mask operation.
type MaskResult struct {
	Entries []Entry
	MaskedCount int
}

// Mask returns a copy of entries with sensitive values obscured.
func Mask(entries []Entry, opts MaskOptions) MaskResult {
	if opts.Placeholder == "" {
		opts.Placeholder = "****"
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	out := make([]Entry, len(entries))
	masked := 0

	for i, e := range entries {
		if shouldMask(e.Key, keySet, opts.Patterns) {
			out[i] = Entry{
				Key:     e.Key,
				Value:   maskValue(e.Value, opts.Placeholder, opts.ShowLast),
				Comment: e.Comment,
			}
			masked++
		} else {
			out[i] = e
		}
	}

	return MaskResult{Entries: out, MaskedCount: masked}
}

// FormatMaskSummary returns a human-readable summary of the mask operation.
func FormatMaskSummary(r MaskResult) string {
	return fmt.Sprintf("masked %d key(s)", r.MaskedCount)
}

func shouldMask(key string, keySet map[string]struct{}, patterns []string) bool {
	if _, ok := keySet[strings.ToUpper(key)]; ok {
		return true
	}
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func maskValue(val, placeholder string, showLast int) string {
	if showLast <= 0 || showLast >= len(val) {
		return placeholder
	}
	return placeholder + val[len(val)-showLast:]
}
