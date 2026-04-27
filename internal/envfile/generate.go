package envfile

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// GenerateOptions controls how placeholder entries are generated.
type GenerateOptions struct {
	// Keys is the list of key names to generate.
	Keys []string
	// DefaultValue is used when Type is "literal".
	DefaultValue string
	// Type is one of: "literal", "random", "uuid".
	Type string
	// RandomLength is the byte length used for random hex generation.
	RandomLength int
	// Overwrite replaces existing keys when true.
	Overwrite bool
}

// GenerateResult holds the outcome of a Generate call.
type GenerateResult struct {
	Entries  []Entry
	Added    []string
	Skipped  []string
}

// Generate creates or replaces entries in base according to opts.
func Generate(base []Entry, opts GenerateOptions) (GenerateResult, error) {
	if opts.Type == "" {
		opts.Type = "literal"
	}
	if opts.RandomLength <= 0 {
		opts.RandomLength = 16
	}

	existing := make(map[string]int, len(base))
	for i, e := range base {
		existing[e.Key] = i
	}

	out := make([]Entry, len(base))
	copy(out, base)

	var added, skipped []string

	for _, key := range opts.Keys {
		val, err := generateValue(opts)
		if err != nil {
			return GenerateResult{}, fmt.Errorf("generate value for %q: %w", key, err)
		}

		if idx, exists := existing[key]; exists {
			if !opts.Overwrite {
				skipped = append(skipped, key)
				continue
			}
			out[idx].Value = val
			added = append(added, key)
		} else {
			out = append(out, Entry{Key: key, Value: val})
			existing[key] = len(out) - 1
			added = append(added, key)
		}
	}

	return GenerateResult{Entries: out, Added: added, Skipped: skipped}, nil
}

func generateValue(opts GenerateOptions) (string, error) {
	switch strings.ToLower(opts.Type) {
	case "literal":
		return opts.DefaultValue, nil
	case "random":
		b := make([]byte, opts.RandomLength)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b), nil
	case "uuid":
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		b[6] = (b[6] & 0x0f) | 0x40
		b[8] = (b[8] & 0x3f) | 0x80
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
	default:
		return "", fmt.Errorf("unknown generate type %q", opts.Type)
	}
}
