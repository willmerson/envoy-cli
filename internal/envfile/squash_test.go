package envfile

import (
	"testing"
)

func squashEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "DB_HOST", Value: "remotehost"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_NAME", Value: "myapp-v2"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestSquash_KeepFirst(t *testing.T) {
	res := Squash(squashEntries(), SquashOptions{})
	if res.Squashed != 2 {
		t.Fatalf("expected 2 squashed, got %d", res.Squashed)
	}
	if res.Kept != 4 {
		t.Fatalf("expected 4 kept, got %d", res.Kept)
	}
	// First occurrence of DB_HOST should be "localhost"
	for _, e := range res.Entries {
		if e.Key == "DB_HOST" && e.Value != "localhost" {
			t.Errorf("expected DB_HOST=localhost, got %s", e.Value)
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME=myapp, got %s", e.Value)
		}
	}
}

func TestSquash_KeepLast(t *testing.T) {
	res := Squash(squashEntries(), SquashOptions{KeepLast: true})
	if res.Squashed != 2 {
		t.Fatalf("expected 2 squashed, got %d", res.Squashed)
	}
	for _, e := range res.Entries {
		if e.Key == "DB_HOST" && e.Value != "remotehost" {
			t.Errorf("expected DB_HOST=remotehost, got %s", e.Value)
		}
		if e.Key == "APP_NAME" && e.Value != "myapp-v2" {
			t.Errorf("expected APP_NAME=myapp-v2, got %s", e.Value)
		}
	}
}

func TestSquash_PrefixGroups(t *testing.T) {
	res := Squash(squashEntries(), SquashOptions{PrefixGroups: true, Separator: "_"})
	if res.Squashed != 2 {
		t.Fatalf("expected 2 squashed, got %d", res.Squashed)
	}
	if res.Kept != 4 {
		t.Fatalf("expected 4 kept, got %d", res.Kept)
	}
}

func TestSquash_NoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "1"},
		{Key: "BAR", Value: "2"},
	}
	res := Squash(entries, SquashOptions{})
	if res.Squashed != 0 {
		t.Errorf("expected 0 squashed, got %d", res.Squashed)
	}
	if res.Kept != 2 {
		t.Errorf("expected 2 kept, got %d", res.Kept)
	}
}

func TestSquash_Empty(t *testing.T) {
	res := Squash([]Entry{}, SquashOptions{})
	if len(res.Entries) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestFormatSquashSummary_WithDuplicates(t *testing.T) {
	res := SquashResult{Kept: 4, Squashed: 2}
	msg := FormatSquashSummary(res)
	expected := "squash: removed 2 duplicate key(s), 4 unique key(s) kept"
	if msg != expected {
		t.Errorf("got %q, want %q", msg, expected)
	}
}

func TestFormatSquashSummary_NoDuplicates(t *testing.T) {
	res := SquashResult{Kept: 3, Squashed: 0}
	msg := FormatSquashSummary(res)
	if msg != "squash: no duplicate keys found" {
		t.Errorf("unexpected message: %s", msg)
	}
}
