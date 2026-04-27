package envfile

import (
	"fmt"
	"strings"
)

// PinResult describes the outcome of a pin operation on a single entry.
type PinResult struct {
	Key     string
	OldValue string
	NewValue string
	Pinned  bool
}

// PinOptions controls how Pin behaves.
type PinOptions struct {
	// Keys is the explicit list of keys to pin. If empty, all keys are pinned.
	Keys []string
	// Prefix restricts pinning to keys that start with this prefix.
	Prefix string
	// Overwrite replaces an already-pinned value with the new one.
	Overwrite bool
}

// Pin locks the values of selected entries to the provided pinMap (key → value).
// Entries not present in pinMap are left untouched.
// Returns the updated entries and a slice of PinResults for reporting.
func Pin(entries []Entry, pinMap map[string]string, opts PinOptions) ([]Entry, []PinResult) {
	wantKey := buildKeySet(opts.Keys)
	results := make([]PinResult, 0)

	out := make([]Entry, len(entries))
	for i, e := range entries {
		out[i] = e

		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if len(wantKey) > 0 && !wantKey[e.Key] {
			continue
		}

		newVal, ok := pinMap[e.Key]
		if !ok {
			continue
		}

		if e.Value == newVal && !opts.Overwrite {
			results = append(results, PinResult{Key: e.Key, OldValue: e.Value, NewValue: newVal, Pinned: false})
			continue
		}

		out[i].Value = newVal
		results = append(results, PinResult{Key: e.Key, OldValue: e.Value, NewValue: newVal, Pinned: true})
	}

	return out, results
}

// FormatPinSummary returns a human-readable summary of pin results.
func FormatPinSummary(results []PinResult) string {
	var sb strings.Builder
	pinned := 0
	for _, r := range results {
		if r.Pinned {
			pinned++
			sb.WriteString(fmt.Sprintf("  pinned  %s  %q → %q\n", r.Key, r.OldValue, r.NewValue))
		} else {
			sb.WriteString(fmt.Sprintf("  skipped %s  (already %q)\n", r.Key, r.OldValue))
		}
	}
	sb.WriteString(fmt.Sprintf("%d key(s) pinned.\n", pinned))
	return sb.String()
}

func buildKeySet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
