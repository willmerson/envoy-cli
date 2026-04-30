package envfile

import (
	"testing"
)

func scopeEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "db.local"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_DEBUG", Value: "true"},
	}
}

func TestScope_ByPrefix(t *testing.T) {
	res := Scope(scopeEntries(), ScopeOptions{Prefix: "APP_", CaseSensitive: true})
	if res.Matched != 3 {
		t.Fatalf("expected 3 matched, got %d", res.Matched)
	}
	for _, e := range res.Entries {
		if !startsWith(e.Key, "APP_") {
			t.Errorf("unexpected key %q in result", e.Key)
		}
	}
}

func TestScope_StripPrefix(t *testing.T) {
	res := Scope(scopeEntries(), ScopeOptions{Prefix: "APP_", StripPrefix: true, CaseSensitive: true})
	expected := []string{"HOST", "PORT", "DEBUG"}
	for i, e := range res.Entries {
		if e.Key != expected[i] {
			t.Errorf("expected key %q, got %q", expected[i], e.Key)
		}
	}
}

func TestScope_CaseInsensitive(t *testing.T) {
	res := Scope(scopeEntries(), ScopeOptions{Prefix: "app_", CaseSensitive: false})
	if res.Matched != 3 {
		t.Fatalf("expected 3 matched, got %d", res.Matched)
	}
}

func TestScope_EmptyPrefix(t *testing.T) {
	res := Scope(scopeEntries(), ScopeOptions{Prefix: ""})
	if res.Matched != res.Total {
		t.Errorf("expected all entries when prefix is empty")
	}
}

func TestScope_NoMatch(t *testing.T) {
	res := Scope(scopeEntries(), ScopeOptions{Prefix: "REDIS_", CaseSensitive: true})
	if res.Matched != 0 {
		t.Errorf("expected 0 matches, got %d", res.Matched)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected empty entries slice")
	}
}

func TestFormatScopeSummary_WithPrefix(t *testing.T) {
	res := ScopeResult{Total: 5, Matched: 3}
	s := FormatScopeSummary(res, "APP_")
	if s != `scoped 3/5 entries matching prefix "APP_"` {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestFormatScopeSummary_NoPrefix(t *testing.T) {
	res := ScopeResult{Total: 4, Matched: 4}
	s := FormatScopeSummary(res, "")
	if s != "scoped all 4 entries" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
