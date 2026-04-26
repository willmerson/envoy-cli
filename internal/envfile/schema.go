package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaFieldType represents the expected type of a .env value.
type SchemaFieldType string

const (
	TypeString  SchemaFieldType = "string"
	TypeInt     SchemaFieldType = "int"
	TypeBool    SchemaFieldType = "bool"
	TypeURL     SchemaFieldType = "url"
	TypeEmail   SchemaFieldType = "email"
	TypePattern SchemaFieldType = "pattern"
)

// SchemaField describes the validation rules for a single env key.
type SchemaField struct {
	Key      string          // env key name
	Type     SchemaFieldType // expected value type
	Required bool            // whether the key must be present
	Pattern  string          // regex pattern (used when Type == TypePattern)
	Allowed  []string        // optional allowlist of accepted values
}

// SchemaViolation describes a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

var (
	reURL   = regexp.MustCompile(`^https?://[^\s]+$`)
	reEmail = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	reInt   = regexp.MustCompile(`^-?\d+$`)
	reBool  = regexp.MustCompile(`^(?i)(true|false|1|0|yes|no)$`)
)

// ValidateSchema checks a slice of Entry values against the provided schema
// fields. It returns a list of violations; an empty slice means the entries
// are fully compliant with the schema.
func ValidateSchema(entries []Entry, fields []SchemaField) []SchemaViolation {
	keyMap := make(map[string]string, len(entries))
	for _, e := range entries {
		keyMap[e.Key] = e.Value
	}

	var violations []SchemaViolation

	for _, field := range fields {
		val, present := keyMap[field.Key]

		if !present {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing",
				})
			}
			continue
		}

		// Allowlist check
		if len(field.Allowed) > 0 {
			if !containsString(field.Allowed, val) {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("value %q is not in allowed set [%s]", val, strings.Join(field.Allowed, ", ")),
				})
			}
			continue
		}

		// Type check
		switch field.Type {
		case TypeInt:
			if !reInt.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected integer, got %q", val)})
			}
		case TypeBool:
			if !reBool.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected boolean, got %q", val)})
			}
		case TypeURL:
			if !reURL.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected URL, got %q", val)})
			}
		case TypeEmail:
			if !reEmail.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected email, got %q", val)})
			}
		case TypePattern:
			if field.Pattern == "" {
				break
			}
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err)})
				break
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern)})
			}
		}
	}

	return violations
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
