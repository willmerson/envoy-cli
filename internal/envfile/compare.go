package envfile

import "fmt"

// CompareResult holds the result of comparing two env file sets.
type CompareResult struct {
	OnlyInA    []Entry
	OnlyInB    []Entry
	Different  []EntryPair
	Identical  []Entry
}

// EntryPair holds two entries with the same key but different values.
type EntryPair struct {
	A Entry
	B Entry
}

// Compare performs a full comparison between two slices of entries.
// It identifies keys only in A, only in B, present in both with different values,
// and present in both with identical values.
func Compare(a, b []Entry) CompareResult {
	mapA := make(map[string]Entry)
	mapB := make(map[string]Entry)

	for _, e := range a {
		mapA[e.Key] = e
	}
	for _, e := range b {
		mapB[e.Key] = e
	}

	result := CompareResult{}

	for _, e := range a {
		if e.Key == "" {
			continue
		}
		if eb, ok := mapB[e.Key]; ok {
			if e.Value == eb.Value {
				result.Identical = append(result.Identical, e)
			} else {
				result.Different = append(result.Different, EntryPair{A: e, B: eb})
			}
		} else {
			result.OnlyInA = append(result.OnlyInA, e)
		}
	}

	for _, e := range b {
		if e.Key == "" {
			continue
		}
		if _, ok := mapA[e.Key]; !ok {
			result.OnlyInB = append(result.OnlyInB, e)
		}
	}

	return result
}

// FormatCompare returns a human-readable summary of a CompareResult.
func FormatCompare(r CompareResult, labelA, labelB string) string {
	var out string

	if len(r.OnlyInA) > 0 {
		out += fmt.Sprintf("Only in %s:\n", labelA)
		for _, e := range r.OnlyInA {
			out += fmt.Sprintf("  + %s=%s\n", e.Key, e.Value)
		}
	}

	if len(r.OnlyInB) > 0 {
		out += fmt.Sprintf("Only in %s:\n", labelB)
		for _, e := range r.OnlyInB {
			out += fmt.Sprintf("  + %s=%s\n", e.Key, e.Value)
		}
	}

	if len(r.Different) > 0 {
		out += "Changed:\n"
		for _, p := range r.Different {
			out += fmt.Sprintf("  ~ %s: %s -> %s\n", p.A.Key, p.A.Value, p.B.Value)
		}
	}

	if len(r.Identical) > 0 {
		out += fmt.Sprintf("Identical: %d key(s)\n", len(r.Identical))
	}

	if out == "" {
		out = "No differences found.\n"
	}

	return out
}
