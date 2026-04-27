package envfile

import (
	"testing"
)

func maskEntries() []Entry {
	return []Entry{
		{Key: "DATABASE_URL", Value: "postgres://user:secret@host/db"},
		{Key: "API_KEY", Value: "sk-abcdef123456"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "SECRET_TOKEN", Value: "topsecret"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestMask_ExactKey(t *testing.T) {
	res := Mask(maskEntries(), MaskOptions{Keys: []string{"API_KEY"}})
	if res.MaskedCount != 1 {
		t.Fatalf("expected 1 masked, got %d", res.MaskedCount)
	}
	for _, e := range res.Entries {
		if e.Key == "API_KEY" && e.Value != "****" {
			t.Errorf("expected **** got %s", e.Value)
		}
	}
}

func TestMask_PatternMatch(t *testing.T) {
	res := Mask(maskEntries(), MaskOptions{Patterns: []string{"SECRET"}})
	// SECRET_TOKEN and DATABASE_URL does not match; only SECRET_TOKEN matches
	if res.MaskedCount != 1 {
		t.Fatalf("expected 1 masked, got %d", res.MaskedCount)
	}
}

func TestMask_ShowLast(t *testing.T) {
	res := Mask(maskEntries(), MaskOptions{
		Keys:     []string{"API_KEY"},
		ShowLast: 4,
	})
	for _, e := range res.Entries {
		if e.Key == "API_KEY" {
			if e.Value != "****3456" {
				t.Errorf("expected ****3456 got %s", e.Value)
			}
		}
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	res := Mask(maskEntries(), MaskOptions{
		Keys:        []string{"PORT"},
		Placeholder: "[REDACTED]",
	})
	for _, e := range res.Entries {
		if e.Key == "PORT" && e.Value != "[REDACTED]" {
			t.Errorf("expected [REDACTED] got %s", e.Value)
		}
	}
}

func TestMask_UnaffectedKeysUnchanged(t *testing.T) {
	res := Mask(maskEntries(), MaskOptions{Keys: []string{"API_KEY"}})
	for _, e := range res.Entries {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("APP_NAME should be unchanged, got %s", e.Value)
		}
	}
}

func TestFormatMaskSummary(t *testing.T) {
	res := MaskResult{MaskedCount: 3}
	s := FormatMaskSummary(res)
	if s != "masked 3 key(s)" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestMask_CaseInsensitiveKey(t *testing.T) {
	res := Mask(maskEntries(), MaskOptions{Keys: []string{"api_key"}})
	if res.MaskedCount != 1 {
		t.Fatalf("expected 1 masked, got %d", res.MaskedCount)
	}
}
