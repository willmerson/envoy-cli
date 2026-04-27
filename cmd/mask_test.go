package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/envfile"
)

func writeTempEnvForMask(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestMaskCommand_ExactKey(t *testing.T) {
	src := writeTempEnvForMask(t, "API_KEY=sk-secret\nAPP_NAME=myapp\n")
	out := filepath.Join(t.TempDir(), "masked.env")

	rootCmd.SetArgs([]string{"mask", "--env", src, "--key", "API_KEY", "--output", out})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := envfile.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Key == "API_KEY" && e.Value != "****" {
			t.Errorf("expected **** got %s", e.Value)
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("APP_NAME should be unchanged")
		}
	}
}

func TestMaskCommand_PatternMask(t *testing.T) {
	src := writeTempEnvForMask(t, "DB_PASSWORD=hunter2\nDB_HOST=localhost\n")
	out := filepath.Join(t.TempDir(), "masked.env")

	rootCmd.SetArgs([]string{"mask", "--env", src, "--pattern", "PASSWORD", "--output", out})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := envfile.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Key == "DB_PASSWORD" && e.Value != "****" {
			t.Errorf("DB_PASSWORD should be masked, got %s", e.Value)
		}
	}
}

func TestMaskCommand_ShowLast(t *testing.T) {
	src := writeTempEnvForMask(t, "TOKEN=abcdef7890\n")
	out := filepath.Join(t.TempDir(), "masked.env")

	rootCmd.SetArgs([]string{"mask", "--env", src, "--key", "TOKEN", "--show-last", "4", "--output", out})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(out)
	if !strings.Contains(string(b), "7890") {
		t.Errorf("expected last 4 chars to be visible in output")
	}
}
