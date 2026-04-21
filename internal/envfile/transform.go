package envfile

import (
	"strings"
)

// TransformOption controls how values are transformed.
type TransformOption int

const (
	TransformUppercase TransformOption = iota
	TransformLowercase
	TransformTrimSpace
	TransformQuoteAll
	TransformUnquoteAll
)

// TransformResult holds the result of a transform operation.
type TransformResult struct {
	Modified int
	Total    int
}

// Transform applies one or more transformations to the values (or keys) of entries.
// If keys is true, transformations are applied to keys instead of values.
func Transform(entries []Entry, keys bool, opts ...TransformOption) ([]Entry, TransformResult) {
	result := make([]Entry, len(entries))
	modified := 0

	for i, e := range entries {
		orig := e
		for _, opt := range opts {
			if keys {
				e.Key = applyTransform(e.Key, opt)
			} else {
				e.Value = applyTransform(e.Value, opt)
			}
		}
		result[i] = e
		if e.Key != orig.Key || e.Value != orig.Value {
			modified++
		}
	}

	return result, TransformResult{Modified: modified, Total: len(entries)}
}

func applyTransform(s string, opt TransformOption) string {
	switch opt {
	case TransformUppercase:
		return strings.ToUpper(s)
	case TransformLowercase:
		return strings.ToLower(s)
	case TransformTrimSpace:
		return strings.TrimSpace(s)
	case TransformQuoteAll:
		if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
			return s
		}
		return `"` + s + `"`
	case TransformUnquoteAll:
		if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		return s
	}
	return s
}
