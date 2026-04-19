package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// EnvFile holds all entries parsed from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
}

// Parse reads and parses a .env file from the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	ef := &EnvFile{Path: path}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			ef.Entries = append(ef.Entries, Entry{Comment: line})
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		ef.Entries = append(ef.Entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}
	return ef, nil
}

// ToMap returns all key-value entries as a map, ignoring comments.
func (ef *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}
