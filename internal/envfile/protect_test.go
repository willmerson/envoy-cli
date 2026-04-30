package envfile

import (
	"strings"
	"testing"
)

func protectEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "API_KEY", Value: "abc123"},
	}
}

func TestProtect_ExactKey(t *testing.T) {
	entries, result, err := Protect(protectEntries(), ProtectOptions{Keys: []string{"DB_PASS"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Protected) != 1 || result.Protected[0] != "DB_PASS" {
		t.Errorf("expected DB_PASS protected, got %v", result.Protected)
	}
	if !IsProtected(entries, "DB_PASS") {
		t.Error("expected DB_PASS to be marked protected")
	}
	if IsProtected(entries, "DB_HOST") {
		t.Error("DB_HOST should not be protected")
	}
}

func TestProtect_PrefixFilter(t *testing.T) {
	_, result, err := Protect(protectEntries(), ProtectOptions{Prefixes: []string{"DB_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Protected) != 2 {
		t.Errorf("expected 2 protected, got %d", len(result.Protected))
	}
}

func TestProtect_DryRun(t *testing.T) {
	entries, result, err := Protect(protectEntries(), ProtectOptions{
		Keys:   []string{"API_KEY"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Protected) != 1 {
		t.Errorf("expected 1 protected in dry-run, got %d", len(result.Protected))
	}
	// Dry-run should not insert sentinel comments
	if IsProtected(entries, "API_KEY") {
		t.Error("dry-run should not insert protection markers")
	}
}

func TestProtect_NoKeysError(t *testing.T) {
	_, _, err := Protect(protectEntries(), ProtectOptions{})
	if err == nil {
		t.Error("expected error when no keys or prefixes given")
	}
}

func TestProtect_AlreadyProtectedNotDuplicated(t *testing.T) {
	base := protectEntries()
	once, _, _ := Protect(base, ProtectOptions{Keys: []string{"DB_HOST"}})
	twice, result, _ := Protect(once, ProtectOptions{Keys: []string{"DB_HOST"}})
	// sentinel should not be doubled
	sentinelCount := 0
	for _, e := range twice {
		if e.Comment && strings.HasPrefix(e.Key, "#PROTECTED: DB_HOST") {
			sentinelCount++
		}
	}
	if sentinelCount != 1 {
		t.Errorf("expected 1 sentinel, got %d", sentinelCount)
	}
	_ = result
}

func TestFormatProtectResult(t *testing.T) {
	r := ProtectResult{Protected: []string{"SECRET_KEY", "DB_PASS"}}
	out := FormatProtectResult(r)
	if !strings.Contains(out, "SECRET_KEY") || !strings.Contains(out, "DB_PASS") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func TestFormatProtectResult_Empty(t *testing.T) {
	out := FormatProtectResult(ProtectResult{})
	if !strings.Contains(out, "No keys protected") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
