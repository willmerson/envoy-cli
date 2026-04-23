package envfile

import (
	"testing"
)

func promoEntries(kvs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestPromote_AddsNewKeys(t *testing.T) {
	src := promoEntries("NEW_KEY", "hello")
	dst := promoEntries("EXISTING", "world")

	result, res, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 1 || res.Updated != 0 || res.Skipped != 0 {
		t.Errorf("unexpected result: %+v", res)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := promoEntries("KEY", "new_value")
	dst := promoEntries("KEY", "old_value")

	result, res, err := Promote(src, dst, PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 || res.Updated != 0 {
		t.Errorf("expected skipped=1, got %+v", res)
	}
	if result[0].Value != "old_value" {
		t.Errorf("expected old_value to be preserved")
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := promoEntries("KEY", "new_value")
	dst := promoEntries("KEY", "old_value")

	result, res, err := Promote(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Updated != 1 {
		t.Errorf("expected updated=1, got %+v", res)
	}
	if result[0].Value != "new_value" {
		t.Errorf("expected new_value, got %s", result[0].Value)
	}
}

func TestPromote_PrefixFilter(t *testing.T) {
	src := promoEntries("PROD_HOST", "prod.example.com", "DEV_HOST", "localhost")
	dst := promoEntries()

	_, res, err := Promote(src, dst, PromoteOptions{PrefixFilter: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 1 || res.Skipped != 1 {
		t.Errorf("expected added=1 skipped=1, got %+v", res)
	}
}

func TestPromote_StripPrefix(t *testing.T) {
	src := promoEntries("STAGING_DB", "staging-db", "STAGING_PORT", "5432")
	dst := promoEntries()

	result, res, err := Promote(src, dst, PromoteOptions{
		PrefixFilter: "STAGING_",
		StripPrefix:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 {
		t.Errorf("expected added=2, got %+v", res)
	}
	if result[0].Key != "DB" || result[1].Key != "PORT" {
		t.Errorf("expected stripped keys, got %v %v", result[0].Key, result[1].Key)
	}
}

func TestFormatPromoteResult(t *testing.T) {
	r := PromoteResult{Added: 3, Updated: 1, Skipped: 2}
	s := FormatPromoteResult(r)
	expected := "promoted: 3 added, 1 updated, 2 skipped"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
