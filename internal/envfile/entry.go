package envfile

// Entry represents a single line in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment bool   // true if this line is a comment or blank
	Raw     string // original raw line, used for comment lines
}

// ToMap converts a slice of entries into a key/value map.
// Comment entries are skipped.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if !e.Comment {
			m[e.Key] = e.Value
		}
	}
	return m
}

// FromMap converts a map into a slice of entries (order not guaranteed).
func FromMap(m map[string]string) []Entry {
	entries := make([]Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	return entries
}
