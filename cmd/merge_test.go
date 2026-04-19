package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envfile"
)

func TestMergeCommand_StrategyOurs(t *testing.T) {
	base := []envfile.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	override := []envfile.Entry{{Key: "B", Value: "99"}, {Key: "C", Value: "3"}}

	r := envfile.Merge(base, override, envfile.StrategyOurs)
	if len(r.Conflicts) != 1 || r.Conflicts[0] != "B" {
		t.Errorf("expected conflict on B, got %v", r.Conflicts)
	}
	m := envfile.ToMap(r.Entries)
	if m["B"] != "2" {
		t.Errorf("expected B=2 with ours strategy, got %s", m["B"])
	}
	if m["C"] != "3" {
		t.Errorf("expected C=3 added, got %s", m["C"])
	}
}

func TestMergeCommand_StrategyTheirs(t *testing.T) {
	base := []envfile.Entry{{Key: "A", Value: "1"}}
	override := []envfile.Entry{{Key: "A", Value: "overridden"}}

	r := envfile.Merge(base, override, envfile.StrategyTheirs)
	m := envfile.ToMap(r.Entries)
	if m["A"] != "overridden" {
		t.Errorf("expected A=overridden, got %s", m["A"])
	}
}

func TestWriteToBuffer(t *testing.T) {
	entries := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	var buf bytes.Buffer
	if err := envfile.WriteTo(&buf, entries); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
