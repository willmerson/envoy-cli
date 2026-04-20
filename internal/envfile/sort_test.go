package envfile

import (
	"testing"
)

func sortEntries() []Entry {
	return []Entry{
		{Key: "ZEBRA", Value: "1"},
		{Key: "APPLE", Value: "2"},
		{Key: "MANGO", Value: "3"},
		{Key: "BANANA", Value: "4"},
	}
}

func TestSort_Ascending(t *testing.T) {
	entries := sortEntries()
	res := Sort(entries, SortAsc)

	expected := []string{"APPLE", "BANANA", "MANGO", "ZEBRA"}
	for i, key := range expected {
		if res.Entries[i].Key != key {
			t.Errorf("pos %d: want %s, got %s", i, key, res.Entries[i].Key)
		}
	}
}

func TestSort_Descending(t *testing.T) {
	entries := sortEntries()
	res := Sort(entries, SortDesc)

	expected := []string{"ZEBRA", "MANGO", "BANANA", "APPLE"}
	for i, key := range expected {
		if res.Entries[i].Key != key {
			t.Errorf("pos %d: want %s, got %s", i, key, res.Entries[i].Key)
		}
	}
}

func TestSort_MovedCount(t *testing.T) {
	entries := sortEntries()
	res := Sort(entries, SortAsc)
	if res.Moved == 0 {
		t.Error("expected at least one entry to have moved")
	}
}

func TestSort_AlreadySorted(t *testing.T) {
	entries := []Entry{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "GAMMA", Value: "3"},
	}
	res := Sort(entries, SortAsc)
	if res.Moved != 0 {
		t.Errorf("expected 0 moved, got %d", res.Moved)
	}
}

func TestSort_PreservesComments(t *testing.T) {
	entries := []Entry{
		{Key: "# comment", Value: ""},
		{Key: "ZEBRA", Value: "1"},
		{Key: "APPLE", Value: "2"},
	}
	res := Sort(entries, SortAsc)
	if res.Entries[0].Key != "# comment" {
		t.Errorf("expected comment to be first, got %s", res.Entries[0].Key)
	}
}

func TestGroupByPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "dev"},
		{Key: "SECRET", Value: "xyz"},
	}
	groups := GroupByPrefix(entries)

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 1 {
		t.Errorf("expected 1 APP entry, got %d", len(groups["APP"]))
	}
	if len(groups[""]) != 1 {
		t.Errorf("expected 1 ungrouped entry, got %d", len(groups[""]))
	}
}
