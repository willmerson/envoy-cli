package envfile

import (
	"regexp"
	"strings"
)

// RedactOptions controls how sensitive values are redacted.
type RedactOptions struct {
	// Keys is an explicit list of key names to redact (case-insensitive).
	Keys []string
	// Patterns is a list of regex patterns matched against key names.
	Patterns []string
	// Placeholder replaces the redacted value. Defaults to "***".
	Placeholder string
}

// RedactResult holds the output entries and metadata.
type RedactResult struct {
	Entries  []Entry
	Redacted int
}

// Redact returns a copy of entries with sensitive values replaced by a placeholder.
// Matching is done against key names using exact matches (case-insensitive) and regex patterns.
func Redact(entries []Entry, opts RedactOptions) RedactResult {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	exactKeys := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		exactKeys[strings.ToUpper(k)] = struct{}{}
	}

	var compiled []*regexp.Regexp
	for _, p := range opts.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}

	result := make([]Entry, len(entries))
	redacted := 0

	for i, e := range entries {
		result[i] = e
		if shouldRedact(e.Key, exactKeys, compiled) {
			result[i].Value = placeholder
			redacted++
		}
	}

	return RedactResult{Entries: result, Redacted: redacted}
}

func shouldRedact(key string, exactKeys map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := exactKeys[strings.ToUpper(key)]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
