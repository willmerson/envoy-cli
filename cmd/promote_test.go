package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/envfile"
)

func writeTempEnvForPromote(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestPromoteCommand_AddsNewKey(t *testing.T) {
	src := writeTempEnvForPromote(t, "NEW_KEY=hello\n")
	dst := writeTempEnvForPromote(t, "EXISTING=world\n")

	rootCmd.SetArgs([]string{"promote", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading dst: %v", err)
	}
	if !strings.Contains(string(data), "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY in dst, got:\n%s", string(data))
	}
}

func TestPromoteCommand_SkipsExistingByDefault(t *testing.T) {
	src := writeTempEnvForPromote(t, "KEY=new_value\n")
	dst := writeTempEnvForPromote(t, "KEY=old_value\n")

	rootCmd.SetArgs([]string{"promote", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := envfile.Parse(readFile(t, dst))
	if len(entries) == 0 || entries[0].Value != "old_value" {
		t.Errorf("expected old_value to be preserved")
	}
}

func TestPromoteCommand_OverwriteFlag(t *testing.T) {
	src := writeTempEnvForPromote(t, "KEY=new_value\n")
	dst := writeTempEnvForPromote(t, "KEY=old_value\n")

	rootCmd.SetArgs([]string{"promote", "--overwrite", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := envfile.Parse(readFile(t, dst))
	if len(entries) == 0 || entries[0].Value != "new_value" {
		t.Errorf("expected new_value after overwrite, got %v", entries)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(data)
}
