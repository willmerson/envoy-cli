package envfile

import (
	"strings"
	"testing"
)

func diffEntries(kvs map[string]string) []Entry {
	var out []Entry
	for k, v := range kvs {
		out = append(out, Entry{Key: k, Value: v})
	}
	return out
}

func TestDiff_Added(t *testing.T) {
	base := diffEntries(map[string]string{"A": "1"})
	other := diffEntries(map[string]string{"A": "1", "B": "2"})
	results := Diff(base, other)
	for _, d := range results {
		if d.Key == "B" && d.Status != DiffAdded {
			t.Errorf("expected B to be added")
		}
	}
}

func TestDiff_Removed(t *testing.T) {
	base := diffEntries(map[string]string{"A": "1", "B": "2"})
	other := diffEntries(map[string]string{"A": "1"})
	results := Diff(base, other)
	for _, d := range results {
		if d.Key == "B" && d.Status != DiffRemoved {
			t.Errorf("expected B to be removed")
		}
	}
}

func TestDiff_Changed(t *testing.T) {
	base := diffEntries(map[string]string{"A": "old"})
	other := diffEntries(map[string]string{"A": "new"})
	results := Diff(base, other)
	if len(results) != 1 || results[0].Status != DiffChanged {
		t.Errorf("expected A to be changed")
	}
}

func TestFormatDiff(t *testing.T) {
	base := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "old"}}
	other := []Entry{{Key: "B", Value: "new"}, {Key: "C", Value: "3"}}
	d := Diff(base, other)
	out := FormatDiff(d)
	if !strings.Contains(out, "- A") {
		t.Errorf("expected removed A in output")
	}
	if !strings.Contains(out, "~ B") {
		t.Errorf("expected changed B in output")
	}
	if !strings.Contains(out, "+ C") {
		t.Errorf("expected added C in output")
	}
}
