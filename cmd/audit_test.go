package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForAudit(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestAuditCommand_Added(t *testing.T) {
	before := writeTempEnvForAudit(t, "FOO=bar\n")
	after := writeTempEnvForAudit(t, "FOO=bar\nNEW_KEY=hello\n")

	rootCmd.SetArgs([]string{"audit", before, after})
	out := captureOutput(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("command failed: %v", err)
		}
	})
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED in output, got: %s", out)
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %s", out)
	}
}

func TestAuditCommand_NoChanges(t *testing.T) {
	path := writeTempEnvForAudit(t, "FOO=bar\nBAZ=qux\n")

	rootCmd.SetArgs([]string{"audit", path, path})
	out := captureOutput(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("command failed: %v", err)
		}
	})
	if !strings.Contains(out, "No changes detected.") {
		t.Errorf("expected 'No changes detected.' in output, got: %s", out)
	}
}

func TestAuditCommand_Changed(t *testing.T) {
	before := writeTempEnvForAudit(t, "FOO=old\n")
	after := writeTempEnvForAudit(t, "FOO=new\n")

	rootCmd.SetArgs([]string{"audit", before, after})
	out := captureOutput(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("command failed: %v", err)
		}
	})
	if !strings.Contains(out, "CHANGED") {
		t.Errorf("expected CHANGED in output, got: %s", out)
	}
}

// captureOutput redirects stdout during fn execution and returns captured string.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	buf := new(strings.Builder)
	tmp := make([]byte, 1024)
	for {
		n, _ := r.Read(tmp)
		if n == 0 {
			break
		}
		buf.Write(tmp[:n])
	}
	return buf.String()
}
