package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umutdz/envoy-cli/internal/envfile"
)

func writeTempPatchEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0644))
	return p
}

func writePatchSpec(t *testing.T, ops []envfile.PatchOp) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "patch.json")
	data, err := json.Marshal(map[string]interface{}{"ops": ops})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(p, data, 0644))
	return p
}

func TestPatchCommand_SetAndDelete(t *testing.T) {
	envP := writeTempPatchEnv(t, "HOST=localhost\nPORT=8080\nDEBUG=true\n")
	ops := []envfile.PatchOp{
		{Op: "set", Key: "PORT", Value: "9090"},
		{Op: "delete", Key: "DEBUG"},
	}
	patchP := writePatchSpec(t, ops)

	envPath = envP
	out := &strings.Builder{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"patch", "--patch-file", patchP})
	err := rootCmd.Execute()
	require.NoError(t, err)

	result, err := envfile.Parse(envP)
	require.NoError(t, err)

	m := envfile.ToMap(result)
	assert.Equal(t, "9090", m["PORT"])
	_, hasDebug := m["DEBUG"]
	assert.False(t, hasDebug)
}

func TestPatchCommand_Rename(t *testing.T) {
	envP := writeTempPatchEnv(t, "HOST=localhost\nPORT=8080\n")
	ops := []envfile.PatchOp{
		{Op: "rename", Key: "HOST", NewKey: "HOSTNAME"},
	}
	patchP := writePatchSpec(t, ops)

	envPath = envP
	rootCmd.SetArgs([]string{"patch", "--patch-file", patchP})
	err := rootCmd.Execute()
	require.NoError(t, err)

	result, err := envfile.Parse(envP)
	require.NoError(t, err)

	m := envfile.ToMap(result)
	assert.Equal(t, "localhost", m["HOSTNAME"])
	_, hasOld := m["HOST"]
	assert.False(t, hasOld)
}
