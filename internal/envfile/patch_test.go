package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func patchEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestPatch_SetExisting(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "set", Key: "PORT", Value: "9090"}}
	out, res := Patch(entries, ops)
	assert.Equal(t, "9090", findValue(out, "PORT"))
	assert.Contains(t, res.Applied, "set PORT")
}

func TestPatch_SetNew(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "hello"}}
	out, res := Patch(entries, ops)
	assert.Equal(t, "hello", findValue(out, "NEW_KEY"))
	assert.Contains(t, res.Applied, "set NEW_KEY")
}

func TestPatch_Delete(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "delete", Key: "DEBUG"}}
	out, res := Patch(entries, ops)
	for _, e := range out {
		assert.NotEqual(t, "DEBUG", e.Key)
	}
	assert.Contains(t, res.Applied, "delete DEBUG")
}

func TestPatch_DeleteNotFound(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "delete", Key: "MISSING"}}
	_, res := Patch(entries, ops)
	assert.Len(t, res.Applied, 0)
	assert.Contains(t, res.Skipped[0], "MISSING")
}

func TestPatch_Rename(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "rename", Key: "HOST", NewKey: "HOSTNAME"}}
	out, res := Patch(entries, ops)
	assert.Equal(t, "localhost", findValue(out, "HOSTNAME"))
	assert.Contains(t, res.Applied, "rename HOST -> HOSTNAME")
}

func TestPatch_RenameEmptyNewKey(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "rename", Key: "HOST", NewKey: ""}}
	_, res := Patch(entries, ops)
	assert.Len(t, res.Errors, 1)
}

func TestPatch_UnknownOp(t *testing.T) {
	entries := patchEntries()
	ops := []PatchOp{{Op: "upsert", Key: "X"}}
	_, res := Patch(entries, ops)
	assert.Len(t, res.Errors, 1)
	assert.Contains(t, res.Errors[0], "unknown op")
}

func findValue(entries []Entry, key string) string {
	for _, e := range entries {
		if e.Key == key {
			return e.Value
		}
	}
	return ""
}
