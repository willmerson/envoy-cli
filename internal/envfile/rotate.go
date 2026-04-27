package envfile

import (
	"fmt"
	"strings"
)

// RotateResult holds the outcome of a key rotation operation.
type RotateResult struct {
	Rotated []string
	Skipped []string
	Total   int
}

// RotateOptions controls how rotation behaves.
type RotateOptions struct {
	// KeyMap maps old key names to new key names.
	KeyMap map[string]string
	// Prefix rotates all keys sharing this prefix by stripping/replacing it.
	OldPrefix string
	NewPrefix string
	// FailOnMissing returns an error if a mapped key is not found.
	FailOnMissing bool
}

// Rotate renames keys in entries according to the provided options.
// KeyMap takes precedence over prefix-based rotation.
func Rotate(entries []Entry, opts RotateOptions) ([]Entry, RotateResult, error) {
	result := RotateResult{Total: len(entries)}
	out := make([]Entry, len(entries))
	copy(out, entries)

	if len(opts.KeyMap) > 0 {
		seen := make(map[string]bool)
		for i, e := range out {
			if newKey, ok := opts.KeyMap[e.Key]; ok {
				out[i].Key = newKey
				result.Rotated = append(result.Rotated, fmt.Sprintf("%s -> %s", e.Key, newKey))
				seen[e.Key] = true
			}
		}
		if opts.FailOnMissing {
			for old := range opts.KeyMap {
				if !seen[old] {
					return nil, result, fmt.Errorf("rotate: key %q not found in entries", old)
				}
			}
		}
		return out, result, nil
	}

	if opts.OldPrefix != "" {
		for i, e := range out {
			if strings.HasPrefix(e.Key, opts.OldPrefix) {
				newKey := opts.NewPrefix + strings.TrimPrefix(e.Key, opts.OldPrefix)
				out[i].Key = newKey
				result.Rotated = append(result.Rotated, fmt.Sprintf("%s -> %s", e.Key, newKey))
			} else {
				result.Skipped = append(result.Skipped, e.Key)
			}
		}
		return out, result, nil
	}

	return out, result, nil
}

// FormatRotateResult returns a human-readable summary of the rotation.
func FormatRotateResult(r RotateResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rotated %d of %d keys\n", len(r.Rotated), r.Total))
	for _, msg := range r.Rotated {
		sb.WriteString(fmt.Sprintf("  ~ %s\n", msg))
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped %d keys\n", len(r.Skipped)))
	}
	return strings.TrimRight(sb.String(), "\n")
}
