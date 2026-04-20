package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envoy-cli/internal/envfile"
)

func writeTempEnvForRename(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestRenameCommand_Success(t *testing.T) {
	p := writeTempEnvForRename(t, "APP_PORT=8080\nDB_HOST=localhost\n")
	envPath = p

	buf := &bytes.Buffer{}
	renameCmd.SetOut(buf)
	renameCmd.SetArgs([]string{"APP_PORT", "SERVER_PORT"})

	if err := renameCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	m := envfile.ToMap(entries)
	if _, ok := m["SERVER_PORT"]; !ok {
		t.Error("expected SERVER_PORT to exist after rename")
	}
	if _, ok := m["APP_PORT"]; ok {
		t.Error("expected APP_PORT to be gone after rename")
	}
}

func TestRenameCommand_KeyNotFound(t *testing.T) {
	p := writeTempEnvForRename(t, "APP_PORT=8080\n")
	envPath = p

	renameCmd.SetArgs([]string{"MISSING", "NEW_KEY"})
	err := renameCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}
