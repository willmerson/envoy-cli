package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func dedupeEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_ENV", Value: "staging"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "DB_HOST", Value: "remotehost"},
		{Key: "DB_HOST", Value: "finalhost"},
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	entries := dedupeEntries()
	resentries, DedupeKeepFirst)

	assert.Len(t, result, 3)
	assert.Equal(t, "production", result[0].Value) // APP_ENV first
	assert.Equal(t, "localhost", result[1].Value)   // DB_HOST first
	assert.Equal(t, "5432", result[2].Value)         // DB_PORT
	assert.Equal(t, 3, summary.Removed)
}

func TestDedupe_KeepLast(t *testing.T) {
	entries := dedupeEntries()
	result, summary := Dedupe(entries, DedupeKeepLast)

	assert.Len(t, result, 3)
	assert.Equal(t, "staging", result[0].Value)    // APP_ENV last
	assert.Equal(t, "finalhost", result[1].Value)  // DB_HOST last
	assert.Equal(t, "5432", result[2].Value)        // DB_PORT unchanged
	assert.Equal(t, 3, summary.Removed)
}

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	result, summary := Dedupe(entries, DedupeKeepFirst)

	assert.Len(t, result, 2)
	assert.Equal(t, 0, summary.Removed)
	assert.Empty(t, summary.Keys)
}

func TestDedupe_Empty(t *testing.T) {
	result, summary := Dedupe([]Entry{}, DedupeKeepFirst)

	assert.Empty(t, result)
	assert.Equal(t, 0, summary.Removed)
}

func TestDedupe_SummaryKeys(t *testing.T) {
	entries := []Entry{
		{Key: "X", Value: "1"},
		{Key: "X", Value: "2"},
		{Key: "Y", Value: "a"},
		{Key: "Y", Value: "b"},
	}
	_, summary := Dedupe(entries, DedupeKeepFirst)

	assert.Equal(t, 2, summary.Removed)
	assert.Contains(t, summary.Keys, "X")
	assert.Contains(t, summary.Keys, "Y")
}
