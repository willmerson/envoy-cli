package envfile

import (
	"testing"
)

func filterEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_SECRET", Value: ""},
		{Key: "LOG_LEVEL", Value: "debug"},
		{Key: "LOG_FILE", Value: ""},
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Prefix: "DB_"})
	if len(result.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(result.Matched))
	}
	if result.Total != 6 {
		t.Fatalf("expected total 6, got %d", result.Total)
	}
}

func TestFilter_BySuffix(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Suffix: "_HOST"})
	if len(result.Matched) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(result.Matched))
	}
	if result.Matched[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", result.Matched[0].Key)
	}
}

func TestFilter_ByContains_KeysOnly(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Contains: "LOG", KeysOnly: true})
	if len(result.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(result.Matched))
	}
}

func TestFilter_ByContains_IncludeValues(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Contains: "debug"})
	if len(result.Matched) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(result.Matched))
	}
	if result.Matched[0].Key != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL, got %s", result.Matched[0].Key)
	}
}

func TestFilter_EmptyOnly(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{EmptyOnly: true})
	if len(result.Matched) != 2 {
		t.Fatalf("expected 2 empty entries, got %d", len(result.Matched))
	}
}

func TestFilter_NoOptions(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{})
	if len(result.Matched) != 6 {
		t.Fatalf("expected all 6 entries, got %d", len(result.Matched))
	}
}

func TestFilter_PrefixAndEmptyOnly(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Prefix: "LOG_", EmptyOnly: true})
	if len(result.Matched) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(result.Matched))
	}
	if result.Matched[0].Key != "LOG_FILE" {
		t.Errorf("expected LOG_FILE, got %s", result.Matched[0].Key)
	}
}
