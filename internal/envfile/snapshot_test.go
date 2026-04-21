package envfile

import (
	"os"
	"testing"
)

func snapEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET", Value: "s3cr3t"},
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	entries := snapEntries()
	if err := SaveSnapshot("test-snap", ".env", entries); err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	rec, err := LoadSnapshot("test-snap")
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if rec.Meta.Name != "test-snap" {
		t.Errorf("expected name 'test-snap', got %q", rec.Meta.Name)
	}
	if rec.Meta.EntryCount != 3 {
		t.Errorf("expected 3 entries, got %d", rec.Meta.EntryCount)
	}
	if len(rec.Entries) != 3 {
		t.Errorf("expected 3 entries in payload, got %d", len(rec.Entries))
	}
	if rec.Entries[0].Key != "APP_ENV" {
		t.Errorf("unexpected first key: %q", rec.Entries[0].Key)
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_, err := LoadSnapshot("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestListSnapshots(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	entries := snapEntries()
	if err := SaveSnapshot("snap-a", ".env", entries); err != nil {
		t.Fatalf("SaveSnapshot snap-a failed: %v", err)
	}
	if err := SaveSnapshot("snap-b", ".env.prod", entries[:1]); err != nil {
		t.Fatalf("SaveSnapshot snap-b failed: %v", err)
	}

	metas, err := ListSnapshots()
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	if len(metas) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(metas))
	}
}

func TestListSnapshots_Empty(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	metas, err := ListSnapshots()
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(metas))
	}
}

func TestSnapshotDir_Created(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir, err := SnapshotDir()
	if err != nil {
		t.Fatalf("SnapshotDir failed: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected snapshot dir to exist at %s", dir)
	}
}
