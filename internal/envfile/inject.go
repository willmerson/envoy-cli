package envfile

import (
	"fmt"
	"os"
	"strings"
)

// InjectResult holds the result of an inject operation.
type InjectResult struct {
	Injected  int
	Skipped   int
	Overwrite int
}

// InjectOptions controls the behavior of Inject.
type InjectOptions struct {
	// Overwrite existing environment variables if true.
	Overwrite bool
	// Prefix filters which keys to inject (empty means all).
	Prefix string
}

// Inject loads the given entries into the current process environment.
// It returns an InjectResult summarising what happened.
func Inject(entries []Entry, opts InjectOptions) (InjectResult, error) {
	var result InjectResult

	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			result.Skipped++
			continue
		}

		_, exists := os.LookupEnv(e.Key)
		if exists && !opts.Overwrite {
			result.Skipped++
			continue
		}

		if err := os.Setenv(e.Key, e.Value); err != nil {
			return result, fmt.Errorf("inject: failed to set %q: %w", e.Key, err)
		}

		if exists {
			result.Overwrite++
		} else {
			result.Injected++
		}
	}

	return result, nil
}

// Snapshot captures the current values of the given keys from the environment.
// Keys not present in the environment are recorded with an empty value.
func Snapshot(keys []string) map[string]string {
	snap := make(map[string]string, len(keys))
	for _, k := range keys {
		snap[k] = os.Getenv(k)
	}
	return snap
}
