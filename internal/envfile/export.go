package envfile

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ExportFormat defines the output format for exporting env entries.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatExport ExportFormat = "export"
)

// Export converts a slice of Entry into the specified format string.
func Export(entries []Entry, format ExportFormat) (string, error) {
	switch format {
	case FormatDotenv:
		return exportDotenv(entries), nil
	case FormatJSON:
		return exportJSON(entries)
	case FormatExport:
		return exportShell(entries), nil
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportDotenv(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Comment != "" {
			fmt.Fprintf(&sb, "# %s\n", e.Comment)
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, quoteIfNeeded(e.Value))
	}
	return sb.String()
}

func exportShell(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%s\n", e.Key, quoteIfNeeded(e.Value))
	}
	return sb.String()
}

func exportJSON(entries []Entry) (string, error) {
	m := ToMap(entries)
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
