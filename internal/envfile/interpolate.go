package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// InterpolateResult holds the outcome of an interpolation pass.
type InterpolateResult struct {
	Entries  []Entry
	Expanded int
	Unresolved []string
}

var interpolateRe = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// Interpolate expands variable references within entry values using the
// values of other entries in the same slice. References may use either
// ${VAR} or $VAR syntax. If strict is true, unresolved references are
// returned as an error.
func Interpolate(entries []Entry, strict bool) (InterpolateResult, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	result := make([]Entry, len(entries))
	unresolvedSet := map[string]struct{}{}
	expanded := 0

	for i, e := range entries {
		newVal, count, missing := expandValue(e.Value, lookup)
		result[i] = Entry{Key: e.Key, Value: newVal, Comment: e.Comment}
		expanded += count
		for _, m := range missing {
			unresolvedSet[m] = struct{}{}
		}
	}

	unresolved := make([]string, 0, len(unresolvedSet))
	for k := range unresolvedSet {
		unresolved = append(unresolved, k)
	}

	if strict && len(unresolved) > 0 {
		return InterpolateResult{}, fmt.Errorf("unresolved references: %s", strings.Join(unresolved, ", "))
	}

	return InterpolateResult{
		Entries:    result,
		Expanded:   expanded,
		Unresolved: unresolved,
	}, nil
}

func expandValue(val string, lookup map[string]string) (string, int, []string) {
	count := 0
	var missing []string

	result := interpolateRe.ReplaceAllStringFunc(val, func(match string) string {
		subs := interpolateRe.FindStringSubmatch(match)
		key := subs[1]
		if key == "" {
			key = subs[2]
		}
		if v, ok := lookup[key]; ok {
			count++
			return v
		}
		missing = append(missing, key)
		return match
	})
	return result, count, missing
}
