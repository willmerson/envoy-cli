package envfile

import (
	"os"
	"testing"
)

func tmplEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_NAME", Value: "mydb"},
	}
}

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	entries := tmplEntries()
	tmpl := "http://${APP_HOST}:${APP_PORT}/${DB_NAME}"
	res, err := RenderTemplate(tmpl, entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "http://localhost:8080/mydb"
	if res.Rendered != expected {
		t.Errorf("expected %q, got %q", expected, res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestRenderTemplate_DollarStyle(t *testing.T) {
	entries := tmplEntries()
	tmpl := "host=$APP_HOST port=$APP_PORT"
	res, err := RenderTemplate(tmpl, entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "host=localhost port=8080" {
		t.Errorf("unexpected rendered: %q", res.Rendered)
	}
}

func TestRenderTemplate_MissingKey_NonStrict(t *testing.T) {
	entries := tmplEntries()
	tmpl := "value=${MISSING_KEY}"
	res, err := RenderTemplate(tmpl, entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING_KEY" {
		t.Errorf("expected [MISSING_KEY], got %v", res.Missing)
	}
	if res.Rendered != tmpl {
		t.Errorf("expected placeholder preserved, got %q", res.Rendered)
	}
}

func TestRenderTemplate_MissingKey_Strict(t *testing.T) {
	entries := tmplEntries()
	tmpl := "value=${MISSING_KEY}"
	_, err := RenderTemplate(tmpl, entries, true)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestRenderTemplate_FallsBackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	res, err := RenderTemplate("${OS_VAR}", []Entry{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "from-os" {
		t.Errorf("expected 'from-os', got %q", res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	res, err := RenderTemplate("plain text", tmplEntries(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "plain text" {
		t.Errorf("expected 'plain text', got %q", res.Rendered)
	}
}
