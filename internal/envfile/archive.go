package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ArchiveEntry represents a single archived version of an env file.
type ArchiveEntry struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

// ArchiveDir returns the directory used to store archives.
func ArchiveDir(base string) string {
	return filepath.Join(base, ".envoy", "archive")
}

// SaveArchive persists a labelled snapshot of entries to the archive directory.
func SaveArchive(base, label string, entries []Entry) (string, error) {
	dir := ArchiveDir(base)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("archive: mkdir: %w", err)
	}

	now := time.Now().UTC()
	id := fmt.Sprintf("%d", now.UnixNano())

	archive := ArchiveEntry{
		ID:        id,
		Label:     label,
		CreatedAt: now,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(archive, "", "  ")
	if err != nil {
		return "", fmt.Errorf("archive: marshal: %w", err)
	}

	path := filepath.Join(dir, id+".json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("archive: write: %w", err)
	}

	return id, nil
}

// LoadArchive loads a specific archive entry by ID.
func LoadArchive(base, id string) (*ArchiveEntry, error) {
	path := filepath.Join(ArchiveDir(base), id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("archive: %q not found", id)
		}
		return nil, fmt.Errorf("archive: read: %w", err)
	}

	var entry ArchiveEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("archive: unmarshal: %w", err)
	}
	return &entry, nil
}

// ListArchives returns all archive entries sorted by creation time (newest first).
func ListArchives(base string) ([]ArchiveEntry, error) {
	dir := ArchiveDir(base)
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("archive: glob: %w", err)
	}

	var archives []ArchiveEntry
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var entry ArchiveEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		archives = append(archives, entry)
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].CreatedAt.After(archives[j].CreatedAt)
	})
	return archives, nil
}
