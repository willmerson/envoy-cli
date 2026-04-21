package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotMeta holds metadata about a saved snapshot.
type SnapshotMeta struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	File      string    `json:"file"`
	EntryCount int      `json:"entry_count"`
}

// SnapshotRecord is the persisted snapshot payload.
type SnapshotRecord struct {
	Meta    SnapshotMeta `json:"meta"`
	Entries []Entry      `json:"entries"`
}

// SnapshotDir returns the directory used to store snapshots.
func SnapshotDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".envoy", "snapshots")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create snapshot directory: %w", err)
	}
	return dir, nil
}

// SaveSnapshot persists a named snapshot of the given entries.
func SaveSnapshot(name string, sourceFile string, entries []Entry) error {
	dir, err := SnapshotDir()
	if err != nil {
		return err
	}
	rec := SnapshotRecord{
		Meta: SnapshotMeta{
			Name:       name,
			CreatedAt:  time.Now().UTC(),
			File:       sourceFile,
			EntryCount: len(entries),
		},
		Entries: entries,
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}
	dest := filepath.Join(dir, name+".json")
	return os.WriteFile(dest, data, 0600)
}

// LoadSnapshot loads a named snapshot and returns its record.
func LoadSnapshot(name string) (*SnapshotRecord, error) {
	dir, err := SnapshotDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot %q not found", name)
		}
		return nil, fmt.Errorf("failed to read snapshot: %w", err)
	}
	var rec SnapshotRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot: %w", err)
	}
	return &rec, nil
}

// ListSnapshots returns metadata for all saved snapshots.
func ListSnapshots() ([]SnapshotMeta, error) {
	dir, err := SnapshotDir()
	if err != nil {
		return nil, err
	}
	glob := filepath.Join(dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var metas []SnapshotMeta
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			continue
		}
		var rec SnapshotRecord
		if err := json.Unmarshal(data, &rec); err != nil {
			continue
		}
		metas = append(metas, rec.Meta)
	}
	return metas, nil
}
