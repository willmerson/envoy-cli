package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadProfile(t *testing.T) {
	base := t.TempDir()
	p := &Profile{
		Name: "staging",
		Entries: []Entry{
			{Key: "DB_HOST", Value: "localhost"},
			{Key: "DB_PORT", Value: "5432"},
		},
	}

	if err := SaveProfile(base, p); err != nil {
		t.Fatalf("SaveProfile: %v", err)
	}

	loaded, err := LoadProfile(base, "staging")
	if err != nil {
		t.Fatalf("LoadProfile: %v", err)
	}

	if loaded.Name != p.Name {
		t.Errorf("expected name %q, got %q", p.Name, loaded.Name)
	}

	got := make(map[string]string)
	for _, e := range loaded.Entries {
		got[e.Key] = e.Value
	}
	for _, e := range p.Entries {
		if got[e.Key] != e.Value {
			t.Errorf("key %q: expected %q, got %q", e.Key, e.Value, got[e.Key])
		}
	}
}

func TestLoadProfile_NotFound(t *testing.T) {
	base := t.TempDir()
	_, err := LoadProfile(base, "ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestListProfiles(t *testing.T) {
	base := t.TempDir()
	for _, name := range []string{"dev", "prod", "staging"} {
		p := &Profile{Name: name, Entries: []Entry{{Key: "ENV", Value: name}}}
		if err := SaveProfile(base, p); err != nil {
			t.Fatalf("SaveProfile %q: %v", name, err)
		}
	}

	names, err := ListProfiles(base)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(names))
	}
}

func TestListProfiles_Empty(t *testing.T) {
	base := t.TempDir()
	names, err := ListProfiles(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(names))
	}
}

func TestProfileDir(t *testing.T) {
	base := "/tmp/myproject"
	expected := filepath.Join(base, ".envoy", "profiles")
	if got := ProfileDir(base); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
	_ = os.RemoveAll(base)
}
