package envfile

import (
	"fmt"
	"strings"
	"unicode"
)

// LintIssue represents a single linting problem found in an env entry.
type LintIssue struct {
	Key     string
	Message string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("[%s] %s", l.Key, l.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

func (r *LintResult) add(key, msg string) {
	r.Issues = append(r.Issues, LintIssue{Key: key, Message: msg})
}

// Lint checks a slice of Entry values for common style and correctness issues.
func Lint(entries []Entry) LintResult {
	result := LintResult{}
	seen := make(map[string]int)

	for _, e := range entries {
		key := e.Key

		// Duplicate key check
		seen[key]++
		if seen[key] == 2 {
			result.add(key, "duplicate key")
		}

		// Empty key
		if key == "" {
			result.add(key, "empty key")
			continue
		}

		// Key should be uppercase
		if key != strings.ToUpper(key) {
			result.add(key, "key is not uppercase")
		}

		// Key must start with a letter or underscore
		if first := rune(key[0]); !unicode.IsLetter(first) && first != '_' {
			result.add(key, "key must start with a letter or underscore")
		}

		// Key should only contain letters, digits, underscores
		for _, ch := range key {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
				result.add(key, "key contains invalid character: "+string(ch))
				break
			}
		}

		// Warn on empty value
		if strings.TrimSpace(e.Value) == "" {
			result.add(key, "value is empty")
		}

		// Warn on unquoted value containing spaces
		if strings.Contains(e.Value, " ") && !strings.HasPrefix(e.Value, "\"") && !strings.HasPrefix(e.Value, "'") {
			result.add(key, "value contains spaces but is not quoted")
		}
	}

	return result
}
