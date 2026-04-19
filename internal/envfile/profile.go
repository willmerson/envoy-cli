package envfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named environment configuration.
type Profile struct {
	Name    string
	Entries []Entry
}

// Entry holds a single key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// ProfileDir returns the directory used to store profile files.
func ProfileDir(base string) string {
	return filepath.Join(base, ".envoy", "profiles")
}

// LoadProfile reads a named profile from the given base directory.
func LoadProfile(base, name string) (*Profile, error) {
	path := filepath.Join(ProfileDir(base), name+".env")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile %q not found", name)
		}
		return nil, err
	}
	defer f.Close()

	pairs, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parsing profile %q: %w", name, err)
	}

	p := &Profile{Name: name}
	for k, v := range pairs {
		p.Entries = append(p.Entries, Entry{Key: k, Value: v})
	}
	return p, nil
}

// SaveProfile writes a profile to the given base directory.
func SaveProfile(base string, p *Profile) error {
	dir := ProfileDir(base)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating profile dir: %w", err)
	}
	path := filepath.Join(dir, p.Name+".env")
	pairs := make(map[string]string, len(p.Entries))
	for _, e := range p.Entries {
		pairs[e.Key] = e.Value
	}
	return WriteTo(path, pairs)
}

// ListProfiles returns the names of all saved profiles in base.
func ListProfiles(base string) ([]string, error) {
	dir := ProfileDir(base)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".env" {
			names = append(names, e.Name()[:len(e.Name())-4])
		}
	}
	return names, nil
}
