package envfile

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) OK() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) String() string {
	if r.OK() {
		return "validation passed"
	}
	var sb strings.Builder
	for _, e := range r.Errors {
		sb.WriteString("  - " + e.Error() + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Validate checks entries against a required set of keys.
// It reports missing keys and keys with empty values.
func Validate(entries []Entry, required []string) ValidationResult {
	result := ValidationResult{}
	keyMap := make(map[string]string, len(entries))
	for _, e := range entries {
		keyMap[e.Key] = e.Value
	}

	for _, req := range required {
		val, exists := keyMap[req]
		if !exists {
			result.Errors = append(result.Errors, ValidationError{Key: req, Message: "missing required key"})
		} else if strings.TrimSpace(val) == "" {
			result.Errors = append(result.Errors, ValidationError{Key: req, Message: "value is empty"})
		}
	}

	return result
}
