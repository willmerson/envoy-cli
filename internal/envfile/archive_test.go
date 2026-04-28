package envfile

import (
	"testing"
)

func archiveEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestSaveAndLoadArchive(t *testing.T) {
	dir := t.TempDir()
	entries := archiveEntries()

	id, err := SaveArchive(dir, "release-v1", entries)
	if err != nil {
		t.Fatalf("SaveArchive: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty ID")
	}

	loaded, err := LoadArchive(dir, id)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if loaded.Label != "release-v1" {
		t.Errorf("label: got %q, want %q", loaded.Label, "release-v1")
	}
	if len(loaded.Entries) != len(entries) {
		t.Errorf("entries count: got %d, want %d", len(loaded.Entries), len(entries))
	}
	if loaded.Entries[0].Key != "APP_NAME" {
		t.Errorf("first key: got %q, want %q", loaded.Entries[0].Key, "APP_NAME")
	}
}

func TestLoadArchive_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadArchive(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing archive")
	}
}

func TestListArchives(t *testing.T) {
	dir := t.TempDir()
	entries := archiveEntries()

	_, err := SaveArchive(dir, "first", entries)
	if err != nil {
		t.Fatalf("SaveArchive first: %v", err)
	}
	_, err = SaveArchive(dir, "second", entries)
	if err != nil {
		t.Fatalf("SaveArchive second: %v", err)
	}

	list, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("ListArchives: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 archives, got %d", len(list))
	}
	// Newest first
	if list[0].Label != "second" {
		t.Errorf("expected newest first, got %q", list[0].Label)
	}
}

func TestListArchives_Empty(t *testing.T) {
	dir := t.TempDir()
	list, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("ListArchives: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestArchiveDir(t *testing.T) {
	dir := t.TempDir()
	got := ArchiveDir(dir)
	if got == "" {
		t.Fatal("expected non-empty archive dir")
	}
}
