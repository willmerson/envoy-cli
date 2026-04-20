package envfile

import (
	"testing"
)

func rEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestRenameKey_Success(t *testing.T) {
	entries, result, err := RenameKey(rEntries(), "APP_PORT", "SERVER_PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Updated {
		t.Error("expected Updated to be true")
	}
	found := false
	for _, e := range entries {
		if e.Key == "SERVER_PORT" && e.Value == "8080" {
			found = true
		}
		if e.Key == "APP_PORT" {
			t.Error("old key APP_PORT should not exist")
		}
	}
	if !found {
		t.Error("new key SERVER_PORT not found in entries")
	}
}

func TestRenameKey_NotFound(t *testing.T) {
	_, _, err := RenameKey(rEntries(), "MISSING_KEY", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameKey_NewKeyExists(t *testing.T) {
	_, _, err := RenameKey(rEntries(), "APP_NAME", "APP_PORT")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestRenameKey_EmptyOldKey(t *testing.T) {
	_, _, err := RenameKey(rEntries(), "", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for empty old key")
	}
}

func TestRenameKey_SameKey(t *testing.T) {
	_, _, err := RenameKey(rEntries(), "APP_NAME", "APP_NAME")
	if err == nil {
		t.Fatal("expected error when old and new keys are identical")
	}
}

func TestRenameKey_PreservesOrder(t *testing.T) {
	entries, _, err := RenameKey(rEntries(), "APP_PORT", "SERVER_PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "APP_NAME" || entries[1].Key != "SERVER_PORT" || entries[2].Key != "DB_HOST" {
		t.Errorf("order not preserved: %v", entries)
	}
}
