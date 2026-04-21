package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateResult holds the output of a template rendering operation.
type TemplateResult struct {
	Rendered string
	Missing  []string
}

// placeholderRe matches ${VAR_NAME} and $VAR_NAME style placeholders.
var placeholderRe = regexp.MustCompile(`\$\{([A-Z0-9_]+)\}|\$([A-Z0-9_]+)`)

// RenderTemplate replaces placeholders in the template string with values
// from the provided entries. Missing keys are collected in the result.
// If strict is true, missing keys cause an error to be returned.
func RenderTemplate(tmpl string, entries []Entry, strict bool) (TemplateResult, error) {
	envMap := ToMap(entries)
	missingSet := map[string]struct{}{}

	rendered := placeholderRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := extractKey(match)
		if val, ok := envMap[key]; ok {
			return val
		}
		// Fall back to OS environment
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		missingSet[key] = struct{}{}
		return match
	})

	missing := make([]string, 0, len(missingSet))
	for k := range missingSet {
		missing = append(missing, k)
	}

	if strict && len(missing) > 0 {
		return TemplateResult{}, fmt.Errorf("template: unresolved placeholders: %s", strings.Join(missing, ", "))
	}

	return TemplateResult{Rendered: rendered, Missing: missing}, nil
}

// extractKey strips ${ } or $ from a placeholder match.
func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
