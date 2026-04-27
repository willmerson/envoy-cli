package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/envfile"
)

func writeTempEnvForRollback(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRollbackCommand_RestoresSnapshot(t *testing.T) {
	envPath = writeTempEnvForRollback(t, "DB_HOST=new-host\nAPP_ENV=production\n")

	// Save a snapshot with old values.
	snap := []envfile.Entry{
		{Key: "DB_HOST", Value: "old-host"},
		{Key: "APP_ENV", Value: "production"},
	}
	if err := envfile.SaveSnapshot("test-snap", snap); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"--env", envPath, "rollback", "test-snap"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "Rollback complete") {
		t.Errorf("expected rollback summary, got: %s", result)
	}

	entries, err := envfile.Parse(envPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Key == "DB_HOST" && e.Value != "old-host" {
			t.Errorf("expected DB_HOST=old-host, got %s", e.Value)
		}
	}
}

func TestRollbackCommand_DryRun(t *testing.T) {
	envPath = writeTempEnvForRollback(t, "DB_HOST=changed\n")

	snap := []envfile.Entry{{Key: "DB_HOST", Value: "original"}}
	if err := envfile.SaveSnapshot("dry-snap", snap); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	rootCmd.SetArgs([]string{"--env", envPath, "rollback", "--dry-run", "dry-snap"})
	rootCmd.SetOut(&buf)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Dry-run") {
		t.Errorf("expected dry-run output, got: %s", buf.String())
	}
}
