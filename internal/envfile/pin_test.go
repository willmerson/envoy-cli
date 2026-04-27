package envfile

import (
	"strings"
	"testing"
)

func pinEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "APP_VERSION", Value: "1.0.0"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "old-secret"},
	}
}

func TestPin_ExplicitKeys(t *testing.T) {
	pinMap := map[string]string{"APP_VERSION": "2.3.1", "SECRET_KEY": "new-secret"}
	out, results := Pin(pinEntries(), pinMap, PinOptions{Keys: []string{"APP_VERSION", "SECRET_KEY"}, Overwrite: true})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Pinned {
			t.Errorf("expected %s to be pinned", r.Key)
		}
	}
	for _, e := range out {
		if e.Key == "APP_VERSION" && e.Value != "2.3.1" {
			t.Errorf("expected APP_VERSION=2.3.1, got %s", e.Value)
		}
		if e.Key == "SECRET_KEY" && e.Value != "new-secret" {
			t.Errorf("expected SECRET_KEY=new-secret, got %s", e.Value)
		}
	}
}

func TestPin_SkipsAlreadyPinned(t *testing.T) {
	pinMap := map[string]string{"DB_PORT": "5432"}
	_, results := Pin(pinEntries(), pinMap, PinOptions{Overwrite: false})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Pinned {
		t.Error("expected entry to be skipped, not pinned")
	}
}

func TestPin_PrefixFilter(t *testing.T) {
	pinMap := map[string]string{"DB_HOST": "db.prod.internal", "DB_PORT": "5433", "APP_ENV": "production"}
	out, results := Pin(pinEntries(), pinMap, PinOptions{Prefix: "DB_", Overwrite: true})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, e := range out {
		if e.Key == "APP_ENV" && e.Value != "development" {
			t.Error("APP_ENV should not have been modified")
		}
	}
}

func TestPin_AllKeys(t *testing.T) {
	pinMap := map[string]string{
		"APP_ENV": "production",
		"APP_VERSION": "3.0.0",
		"DB_HOST": "prod-db",
		"DB_PORT": "5433",
		"SECRET_KEY": "s3cr3t",
	}
	out, results := Pin(pinEntries(), pinMap, PinOptions{Overwrite: true})

	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
	if len(out) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(out))
	}
}

func TestFormatPinSummary(t *testing.T) {
	results := []PinResult{
		{Key: "APP_VERSION", OldValue: "1.0.0", NewValue: "2.0.0", Pinned: true},
		{Key: "DB_PORT", OldValue: "5432", NewValue: "5432", Pinned: false},
	}
	summary := FormatPinSummary(results)
	if !strings.Contains(summary, "pinned") {
		t.Error("summary should mention pinned")
	}
	if !strings.Contains(summary, "1 key(s) pinned") {
		t.Errorf("unexpected summary: %s", summary)
	}
}
