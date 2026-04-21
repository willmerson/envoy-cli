package envfile

import (
	"os"
	"testing"
)

func injectEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
}

func TestInject_Basic(t *testing.T) {
	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_URL")

	result, err := Inject(injectEntries(), InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Injected != 3 {
		t.Errorf("expected 3 injected, got %d", result.Injected)
	}
	if os.Getenv("APP_HOST") != "localhost" {
		t.Errorf("APP_HOST not set correctly")
	}
}

func TestInject_NoOverwrite(t *testing.T) {
	os.Setenv("APP_HOST", "original")
	t.Cleanup(func() { os.Unsetenv("APP_HOST") })

	entries := []Entry{{Key: "APP_HOST", Value: "new"}}
	result, err := Inject(entries, InjectOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
	if os.Getenv("APP_HOST") != "original" {
		t.Errorf("expected original value to be preserved")
	}
}

func TestInject_Overwrite(t *testing.T) {
	os.Setenv("APP_PORT", "9999")
	t.Cleanup(func() { os.Unsetenv("APP_PORT") })

	entries := []Entry{{Key: "APP_PORT", Value: "8080"}}
	result, err := Inject(entries, InjectOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Overwrite != 1 {
		t.Errorf("expected 1 overwrite, got %d", result.Overwrite)
	}
	if os.Getenv("APP_PORT") != "8080" {
		t.Errorf("expected overwritten value")
	}
}

func TestInject_PrefixFilter(t *testing.T) {
	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_URL")
	t.Cleanup(func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_URL")
	})

	result, err := Inject(injectEntries(), InjectOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Injected != 2 {
		t.Errorf("expected 2 injected, got %d", result.Injected)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
	if os.Getenv("DB_URL") != "" {
		t.Errorf("DB_URL should not have been injected")
	}
}

func TestSnapshot(t *testing.T) {
	os.Setenv("SNAP_A", "alpha")
	os.Unsetenv("SNAP_B")
	t.Cleanup(func() { os.Unsetenv("SNAP_A") })

	snap := Snapshot([]string{"SNAP_A", "SNAP_B"})
	if snap["SNAP_A"] != "alpha" {
		t.Errorf("expected SNAP_A=alpha, got %q", snap["SNAP_A"])
	}
	if snap["SNAP_B"] != "" {
		t.Errorf("expected SNAP_B empty, got %q", snap["SNAP_B"])
	}
}
