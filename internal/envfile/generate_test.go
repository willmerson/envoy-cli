package envfile

import (
	"strings"
	"testing"
)

func genEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestGenerate_LiteralNew(t *testing.T) {
	res, err := Generate(genEntries(), GenerateOptions{
		Keys:         []string{"NEW_KEY"},
		Type:         "literal",
		DefaultValue: "placeholder",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY added, got %v", res.Added)
	}
	m := ToMap(res.Entries)
	if m["NEW_KEY"] != "placeholder" {
		t.Errorf("unexpected value %q", m["NEW_KEY"])
	}
}

func TestGenerate_SkipsExistingWithoutOverwrite(t *testing.T) {
	res, err := Generate(genEntries(), GenerateOptions{
		Keys:         []string{"APP_NAME"},
		Type:         "literal",
		DefaultValue: "changed",
		Overwrite:    false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME skipped, got %v", res.Skipped)
	}
	m := ToMap(res.Entries)
	if m["APP_NAME"] != "myapp" {
		t.Errorf("value should not change, got %q", m["APP_NAME"])
	}
}

func TestGenerate_OverwriteExisting(t *testing.T) {
	res, err := Generate(genEntries(), GenerateOptions{
		Keys:         []string{"PORT"},
		Type:         "literal",
		DefaultValue: "9090",
		Overwrite:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(res.Added))
	}
	m := ToMap(res.Entries)
	if m["PORT"] != "9090" {
		t.Errorf("expected 9090, got %q", m["PORT"])
	}
}

func TestGenerate_RandomValue(t *testing.T) {
	res, err := Generate(genEntries(), GenerateOptions{
		Keys:         []string{"SECRET"},
		Type:         "random",
		RandomLength: 8,
	})
	if err != nil {
		t.Fatal(err)
	}
	m := ToMap(res.Entries)
	if len(m["SECRET"]) != 16 { // 8 bytes => 16 hex chars
		t.Errorf("expected 16 hex chars, got %q", m["SECRET"])
	}
}

func TestGenerate_UUIDValue(t *testing.T) {
	res, err := Generate(genEntries(), GenerateOptions{
		Keys: []string{"REQUEST_ID"},
		Type: "uuid",
	})
	if err != nil {
		t.Fatal(err)
	}
	m := ToMap(res.Entries)
	parts := strings.Split(m["REQUEST_ID"], "-")
	if len(parts) != 5 {
		t.Errorf("expected UUID with 5 parts, got %q", m["REQUEST_ID"])
	}
}

func TestGenerate_UnknownType(t *testing.T) {
	_, err := Generate(genEntries(), GenerateOptions{
		Keys: []string{"X"},
		Type: "bogus",
	})
	if err == nil {
		t.Error("expected error for unknown type")
	}
}
