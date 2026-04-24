package envfile

import (
	"fmt"
	"os"
	"strings"
)

// Write serialises an EnvFile back to disk at ef.Path.
func Write(ef *EnvFile) error {
	return WriteTo(ef, ef.Path)
}

// WriteTo serialises an EnvFile to an arbitrary path.
func WriteTo(ef *EnvFile, path string) error {
	var sb strings.Builder
	for _, e := range ef.Entries {
		if e.Comment != "" {
			sb.WriteString(e.Comment + "\n")
			continue
		}
		value := e.Value
		if strings.ContainsAny(value, " \t") {
			value = fmt.Sprintf("%q", value)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, value))
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("writing env file %q: %w", path, err)
	}
	return nil
}

// Set adds or updates a key in the EnvFile.
func (ef *EnvFile) Set(key, value string) {
	for i, e := range ef.Entries {
		if e.Key == key {
			ef.Entries[i].Value = value
			return
		}
	}
	ef.Entries = append(ef.Entries, Entry{Key: key, Value: value})
}

// Delete removes a key from the EnvFile. Returns true if the key existed.
func (ef *EnvFile) Delete(key string) bool {
	for i, e := range ef.Entries {
		if e.Key == key {
			ef.Entries = append(ef.Entries[:i], ef.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// Get returns the value for the given key and whether the key was found.
func (ef *EnvFile) Get(key string) (string, bool) {
	for _, e := range ef.Entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}
