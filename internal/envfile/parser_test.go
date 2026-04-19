package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	ef, err := envfile.Parse(path)
	require.NoError(t, err)
	m := ef.ToMap()
	assert.Equal(t, "production", m["APP_ENV"])
	assert.Equal(t, "8080", m["PORT"])
}

func TestParse_QuotedValue(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"`+"\n")
	ef, err := envfile.Parse(path)
	require.NoError(t, err)
	assert.Equal(t, "postgres://localhost/mydb", ef.ToMap()["DB_URL"])
}

func TestParse_Comments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nKEY=val\n")
	ef, err := envfile.Parse(path)
	require.NoError(t, err)
	assert.Len(t, ef.Entries, 2)
	assert.Equal(t, "# comment", ef.Entries[0].Comment)
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := envfile.Parse(path)
	assert.Error(t, err)
}

func TestSet_NewKey(t *testing.T) {
	path := writeTempEnv(t, "EXISTING=yes\n")
	ef, err := envfile.Parse(path)
	require.NoError(t, err)
	ef.Set("NEW_KEY", "hello")
	assert.Equal(t, "hello", ef.ToMap()["NEW_KEY"])
}

func TestDelete_ExistingKey(t *testing.T) {
	path := writeTempEnv(t, "TO_DELETE=bye\n")
	ef, err := envfile.Parse(path)
	require.NoError(t, err)
	ok := ef.Delete("TO_DELETE")
	assert.True(t, ok)
	_, exists := ef.ToMap()["TO_DELETE"]
	assert.False(t, exists)
}
