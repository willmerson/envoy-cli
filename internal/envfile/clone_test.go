package envfile

import (
	"testing"
)

func cloneEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestClone_AllEntries(t *testing.T) {
	src := cloneEntries("A", "1", "B", "2")
	dst := cloneEntries("C", "3")

	out, res, err := Clone(src, dst, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 2 || res.Skipped != 0 {
		t.Errorf("expected 2 copied, 0 skipped; got %d/%d", res.Copied, res.Skipped)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := cloneEntries("A", "new")
	dst := cloneEntries("A", "old")

	out, res, err := Clone(src, dst, CloneOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 || res.Copied != 0 {
		t.Errorf("expected 1 skipped; got copied=%d skipped=%d", res.Copied, res.Skipped)
	}
	if out[0].Value != "old" {
		t.Errorf("expected value to remain 'old', got %q", out[0].Value)
	}
}

func TestClone_OverwriteExisting(t *testing.T) {
	src := cloneEntries("A", "new")
	dst := cloneEntries("A", "old")

	out, res, err := Clone(src, dst, CloneOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 1 {
		t.Errorf("expected 1 copied, got %d", res.Copied)
	}
	if out[0].Value != "new" {
		t.Errorf("expected value 'new', got %q", out[0].Value)
	}
}

func TestClone_PrefixFilter(t *testing.T) {
	src := cloneEntries("APP_HOST", "localhost", "APP_PORT", "8080", "DB_URL", "postgres")
	dst := []Entry{}

	out, res, err := Clone(src, dst, CloneOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 2 {
		t.Errorf("expected 2 copied, got %d", res.Copied)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestClone_PrefixRemap(t *testing.T) {
	src := cloneEntries("APP_HOST", "localhost", "APP_PORT", "8080")
	dst := []Entry{}

	out, res, err := Clone(src, dst, CloneOptions{Prefix: "APP_", DestPrefix: "SVC_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 2 {
		t.Errorf("expected 2 copied, got %d", res.Copied)
	}
	if out[0].Key != "SVC_HOST" || out[1].Key != "SVC_PORT" {
		t.Errorf("unexpected keys: %v, %v", out[0].Key, out[1].Key)
	}
}

func TestClone_Empty(t *testing.T) {
	out, res, err := Clone([]Entry{}, []Entry{}, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 0 || len(out) != 0 {
		t.Errorf("expected empty result")
	}
}
