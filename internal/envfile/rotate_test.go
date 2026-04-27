package envfile

import (
	"testing"
)

func rotateEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "db.local"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestRotate_KeyMap(t *testing.T) {
	entries := rotateEntries()
	opts := RotateOptions{
		KeyMap: map[string]string{
			"APP_HOST": "SERVICE_HOST",
			"APP_PORT": "SERVICE_PORT",
		},
	}
	out, res, err := Rotate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if out[0].Key != "SERVICE_HOST" {
		t.Errorf("expected SERVICE_HOST, got %s", out[0].Key)
	}
	if out[1].Key != "SERVICE_PORT" {
		t.Errorf("expected SERVICE_PORT, got %s", out[1].Key)
	}
}

func TestRotate_PrefixSwap(t *testing.T) {
	entries := rotateEntries()
	opts := RotateOptions{
		OldPrefix: "APP_",
		NewPrefix: "SVC_",
	}
	out, res, err := Rotate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if out[0].Key != "SVC_HOST" {
		t.Errorf("expected SVC_HOST, got %s", out[0].Key)
	}
	if len(res.Skipped) != 3 {
		t.Errorf("expected 3 skipped, got %d", len(res.Skipped))
	}
}

func TestRotate_FailOnMissing(t *testing.T) {
	entries := rotateEntries()
	opts := RotateOptions{
		KeyMap:        map[string]string{"NONEXISTENT": "NEW_KEY"},
		FailOnMissing: true,
	}
	_, _, err := Rotate(entries, opts)
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestRotate_NoOp(t *testing.T) {
	entries := rotateEntries()
	out, res, err := Rotate(entries, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 0 {
		t.Errorf("expected 0 rotated, got %d", len(res.Rotated))
	}
	if len(out) != len(entries) {
		t.Errorf("expected same length, got %d", len(out))
	}
}

func TestFormatRotateResult(t *testing.T) {
	r := RotateResult{
		Rotated: []string{"APP_HOST -> SVC_HOST", "APP_PORT -> SVC_PORT"},
		Skipped: []string{"DB_HOST"},
		Total:   5,
	}
	out := FormatRotateResult(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
