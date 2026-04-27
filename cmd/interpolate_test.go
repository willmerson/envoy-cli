package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envoy-cli/internal/envfile"
)

func writeTempEnvForInterpolate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInterpolateCommand_ExpandsRefs(t *testing.T) {
	src := writeTempEnvForInterpolate(t, "BASE=https://example.com\nURL=${BASE}/path\n")
	out := filepath.Join(t.TempDir(), "out.env")

	rootCmd.SetArgs([]string{"--env", src, "interpolate", "--output", out})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	entries, err := envfile.Parse(out)
	if err != nil {
		t.Fatalf("parse output: %v", err)
	}
	for _, e := range entries {
		if e.Key == "URL" && e.Value != "https://example.com/path" {
			t.Errorf("URL: got %q, want %q", e.Value, "https://example.com/path")
		}
	}
}

func TestInterpolateCommand_StrictFails(t *testing.T) {
	src := writeTempEnvForInterpolate(t, "FOO=${UNDEFINED}\n")

	rootCmd.SetArgs([]string{"--env", src, "interpolate", "--strict"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error in strict mode with unresolved ref")
	}
	if !strings.Contains(err.Error(), "unresolved") {
		t.Errorf("error should mention 'unresolved', got: %v", err)
	}
}
