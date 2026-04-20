package envfile

import (
	"strings"
	"testing"
)

var exportEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DEBUG", Value: "true"},
	{Key: "DB_URL", Value: "postgres://localhost/db"},
	{Key: "GREETING", Value: "hello world", Comment: "has a space"},
}

func TestExport_Dotenv(t *testing.T) {
	out, err := Export(exportEntries, FormatDotenv)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", out)
	}
	if !strings.Contains(out, "# has a space") {
		t.Errorf("expected comment in output, got:\n%s", out)
	}
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value for GREETING, got:\n%s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	out, err := Export(exportEntries, FormatExport)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "export APP_NAME=myapp") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	out, err := Export(exportEntries, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"APP_NAME"`) {
		t.Errorf("expected JSON key APP_NAME, got:\n%s", out)
	}
	if !strings.Contains(out, `"myapp"`) {
		t.Errorf("expected JSON value myapp, got:\n%s", out)
	}
}

func TestExport_InvalidFormat(t *testing.T) {
	_, err := Export(exportEntries, "xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExport_Empty(t *testing.T) {
	formats := []Format{FormatDotenv, FormatExport, FormatJSON}
	for _, fmt := range formats {
		out, err := Export([]Entry{}, fmt)
		if err != nil {
			t.Errorf("format %q: unexpected error for empty entries: %v", fmt, err)
		}
		if out == "" {
			continue // empty output is acceptable
		}
		_ = out
	}
}
