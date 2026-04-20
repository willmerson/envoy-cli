package envfile

import (
	"testing"
)

func copyEntries() []Entry {
	return []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/dev"},
		{Key: "API_KEY", Value: "secret123"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestCopyKey_Success(t *testing.T) {
	entries := copyEntries()
	result, info, err := CopyKey(entries, "DATABASE_URL", "DATABASE_URL_BACKUP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
	if info.Key != "DATABASE_URL" || info.NewKey != "DATABASE_URL_BACKUP" {
		t.Errorf("unexpected copy result: %+v", info)
	}
	last := result[len(result)-1]
	if last.Key != "DATABASE_URL_BACKUP" || last.Value != "postgres://localhost/dev" {
		t.Errorf("copied entry mismatch: %+v", last)
	}
}

func TestCopyKey_NotFound(t *testing.T) {
	entries := copyEntries()
	_, _, err := CopyKey(entries, "MISSING_KEY", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestCopyKey_NewKeyExists(t *testing.T) {
	entries := copyEntries()
	_, _, err := CopyKey(entries, "DATABASE_URL", "API_KEY")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestCopyKey_SameKey(t *testing.T) {
	entries := copyEntries()
	_, _, err := CopyKey(entries, "DEBUG", "DEBUG")
	if err == nil {
		t.Fatal("expected error when old and new keys are the same")
	}
}

func TestCopyKey_EmptyOldKey(t *testing.T) {
	entries := copyEntries()
	_, _, err := CopyKey(entries, "", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for empty old key")
	}
}

func TestCopyKey_EmptyNewKey(t *testing.T) {
	entries := copyEntries()
	_, _, err := CopyKey(entries, "DEBUG", "")
	if err == nil {
		t.Fatal("expected error for empty new key")
	}
}

func TestCopyKey_CommentSet(t *testing.T) {
	entries := copyEntries()
	result, _, err := CopyKey(entries, "API_KEY", "API_KEY_COPY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	last := result[len(result)-1]
	if last.Comment == "" {
		t.Error("expected comment to be set on copied entry")
	}
}
