package envfile

import (
	"strings"
	"testing"
)

func cmpEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCompare_OnlyInA(t *testing.T) {
	a := cmpEntries("FOO", "bar", "ONLY_A", "yes")
	b := cmpEntries("FOO", "bar")
	r := Compare(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Key != "ONLY_A" {
		t.Errorf("expected ONLY_A in OnlyInA, got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 0 {
		t.Errorf("expected no OnlyInB entries")
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := cmpEntries("FOO", "bar")
	b := cmpEntries("FOO", "bar", "ONLY_B", "yes")
	r := Compare(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0].Key != "ONLY_B" {
		t.Errorf("expected ONLY_B in OnlyInB, got %v", r.OnlyInB)
	}
}

func TestCompare_Different(t *testing.T) {
	a := cmpEntries("KEY", "old")
	b := cmpEntries("KEY", "new")
	r := Compare(a, b)
	if len(r.Different) != 1 {
		t.Fatalf("expected 1 different entry, got %d", len(r.Different))
	}
	if r.Different[0].A.Value != "old" || r.Different[0].B.Value != "new" {
		t.Errorf("unexpected values in Different: %v", r.Different[0])
	}
}

func TestCompare_Identical(t *testing.T) {
	a := cmpEntries("KEY", "same")
	b := cmpEntries("KEY", "same")
	r := Compare(a, b)
	if len(r.Identical) != 1 || r.Identical[0].Key != "KEY" {
		t.Errorf("expected KEY in Identical, got %v", r.Identical)
	}
	if len(r.Different) != 0 {
		t.Errorf("expected no Different entries")
	}
}

func TestFormatCompare_NoDifferences(t *testing.T) {
	a := cmpEntries("K", "v")
	b := cmpEntries("K", "v")
	r := Compare(a, b)
	out := FormatCompare(r, "A", "B")
	if !strings.Contains(out, "No differences") && !strings.Contains(out, "Identical") {
		t.Errorf("expected no-differences message, got: %s", out)
	}
}

func TestFormatCompare_ShowsLabels(t *testing.T) {
	a := cmpEntries("ONLY_A", "1")
	b := cmpEntries("ONLY_B", "2")
	r := Compare(a, b)
	out := FormatCompare(r, "staging", "production")
	if !strings.Contains(out, "staging") {
		t.Errorf("expected label 'staging' in output, got: %s", out)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected label 'production' in output, got: %s", out)
	}
}
