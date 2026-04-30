package envfile

import (
	"fmt"
	"strings"
)

// ProtectOptions configures which keys are protected from modification.
type ProtectOptions struct {
	Keys     []string // exact key names
	Prefixes []string // key prefixes to protect
	DryRun   bool
}

// ProtectResult holds the outcome of a protect operation.
type ProtectResult struct {
	Protected []string
	Skipped   []string
}

// Protect marks entries as read-only by prepending a "#PROTECTED:" comment
// sentinel above each matched key. Returns entries with protection markers
// inserted and a summary of what was affected.
func Protect(entries []Entry, opts ProtectOptions) ([]Entry, ProtectResult, error) {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		return nil, ProtectResult{}, fmt.Errorf("protect: at least one key or prefix must be specified")
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.TrimSpace(k)] = true
	}

	var result ProtectResult
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if isProtectedEntry(e) {
			out = append(out, e)
			continue
		}
		if shouldProtect(e.Key, keySet, opts.Prefixes) {
			result.Protected = append(result.Protected, e.Key)
			if !opts.DryRun {
				out = append(out, Entry{Key: "#PROTECTED: " + e.Key, Value: "", Comment: true})
			}
			out = append(out, e)
		} else {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
		}
	}

	return out, result, nil
}

// IsProtected reports whether the given key has a protection marker in entries.
func IsProtected(entries []Entry, key string) bool {
	for i, e := range entries {
		if e.Key == key && i > 0 {
			prev := entries[i-1]
			if prev.Comment && strings.HasPrefix(prev.Key, "#PROTECTED: ") {
				return true
			}
		}
	}
	return false
}

// FormatProtectResult returns a human-readable summary.
func FormatProtectResult(r ProtectResult) string {
	var sb strings.Builder
	if len(r.Protected) == 0 {
		sb.WriteString("No keys protected.\n")
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("Protected %d key(s):\n", len(r.Protected)))
	for _, k := range r.Protected {
		sb.WriteString(fmt.Sprintf("  + %s\n", k))
	}
	return sb.String()
}

func shouldProtect(key string, keySet map[string]bool, prefixes []string) bool {
	if keySet[key] {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func isProtectedEntry(e Entry) bool {
	return e.Comment && strings.HasPrefix(e.Key, "#PROTECTED: ")
}
