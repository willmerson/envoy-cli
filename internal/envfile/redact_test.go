package envfile

import (
	"testing"
)

func redactEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "AWS_SECRET_ACCESS_KEY", Value: "aws-secret"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestRedact_ExactKey(t *testing.T) {
	entries := redactEntries()
	res := Redact(entries, RedactOptions{Keys: []string{"DB_PASSWORD"}})
	if res.Redacted != 1 {
		t.Fatalf("expected 1 redacted, got %d", res.Redacted)
	}
	for _, e := range res.Entries {
		if e.Key == "DB_PASSWORD" && e.Value != "***" {
			t.Errorf("expected DB_PASSWORD to be redacted")
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME to be unchanged")
		}
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	entries := redactEntries()
	res := Redact(entries, RedactOptions{Keys: []string{"api_key"}})
	if res.Redacted != 1 {
		t.Fatalf("expected 1 redacted, got %d", res.Redacted)
	}
}

func TestRedact_PatternMatch(t *testing.T) {
	entries := redactEntries()
	res := Redact(entries, RedactOptions{Patterns: []string{"(?i)secret|password"}})
	if res.Redacted != 2 {
		t.Fatalf("expected 2 redacted, got %d", res.Redacted)
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	entries := redactEntries()
	res := Redact(entries, RedactOptions{Keys: []string{"API_KEY"}, Placeholder: "[REDACTED]"})
	for _, e := range res.Entries {
		if e.Key == "API_KEY" && e.Value != "[REDACTED]" {
			t.Errorf("expected custom placeholder, got %s", e.Value)
		}
	}
}

func TestRedact_NoMatch(t *testing.T) {
	entries := redactEntries()
	res := Redact(entries, RedactOptions{Keys: []string{"NONEXISTENT"}})
	if res.Redacted != 0 {
		t.Errorf("expected 0 redacted, got %d", res.Redacted)
	}
	for i, e := range res.Entries {
		if e.Value != entries[i].Value {
			t.Errorf("expected entry %s to be unchanged", e.Key)
		}
	}
}

func TestRedact_OriginalUnmodified(t *testing.T) {
	entries := redactEntries()
	Redact(entries, RedactOptions{Keys: []string{"DB_PASSWORD"}})
	for _, e := range entries {
		if e.Key == "DB_PASSWORD" && e.Value != "supersecret" {
			t.Errorf("original entries should not be modified")
		}
	}
}
