package envfile

import (
	"strings"
)

// TrimOptions controls which trimming operations are applied.
type TrimOptions struct {
	TrimKeySpace   bool // remove leading/trailing whitespace from keys
	TrimValueSpace bool // remove leading/trailing whitespace from values
	TrimQuotes     bool // strip surrounding quotes from values
	TrimPrefix     string // remove a fixed prefix from keys
	TrimSuffix     string // remove a fixed suffix from keys
}

// TrimResult holds the outcome of a Trim operation.
type TrimResult struct {
	Entries  []Entry
	Modified int
}

// Trim applies the given TrimOptions to a slice of entries and returns
// a new slice with the transformations applied along with a summary.
func Trim(entries []Entry, opts TrimOptions) TrimResult {
	out := make([]Entry, len(entries))
	modified := 0

	for i, e := range entries {
		origKey := e.Key
		origVal := e.Value

		if opts.TrimKeySpace {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimPrefix != "" {
			e.Key = strings.TrimPrefix(e.Key, opts.TrimPrefix)
		}
		if opts.TrimSuffix != "" {
			e.Key = strings.TrimSuffix(e.Key, opts.TrimSuffix)
		}
		if opts.TrimValueSpace {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.TrimQuotes {
			e.Value = trimSurroundingQuotes(e.Value)
		}

		if e.Key != origKey || e.Value != origVal {
			modified++
		}
		out[i] = e
	}

	return TrimResult{Entries: out, Modified: modified}
}

// trimSurroundingQuotes removes a matching pair of single or double quotes
// from the start and end of s, if present.
func trimSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// FormatTrimSummary returns a human-readable summary of a TrimResult.
func FormatTrimSummary(r TrimResult) string {
	if r.Modified == 0 {
		return "trim: no entries changed"
	}
	if r.Modified == 1 {
		return "trim: 1 entry changed"
	}
	return strings.Join([]string{
		"trim:",
		strconv(r.Modified),
		"entries changed",
	}, " ")
}

func strconv(n int) string {
	return strings.TrimSpace(strings.Replace(" "+itoa(n), " ", "", -1))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
