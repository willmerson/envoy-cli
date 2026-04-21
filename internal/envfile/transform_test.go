package envfile

import (
	"testing"
)

func tEntries() []Entry {
	return []Entry{
		{Key: "db_host", Value: "localhost"},
		{Key: "db_port", Value: "  5432  "},
		{Key: "APP_NAME", Value: "MyApp"},
		{Key: "debug", Value: "\"true\""},
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	out, res := Transform(tEntries(), true, TransformUppercase)
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out[0].Key)
	}
	if res.Modified != 3 { // APP_NAME already uppercase
		t.Errorf("expected 3 modified, got %d", res.Modified)
	}
}

func TestTransform_LowercaseValues(t *testing.T) {
	out, res := Transform(tEntries(), false, TransformLowercase)
	if out[2].Value != "myapp" {
		t.Errorf("expected myapp, got %s", out[2].Value)
	}
	if res.Total != 4 {
		t.Errorf("expected total 4, got %d", res.Total)
	}
}

func TestTransform_TrimSpaceValues(t *testing.T) {
	out, _ := Transform(tEntries(), false, TransformTrimSpace)
	if out[1].Value != "5432" {
		t.Errorf("expected '5432', got '%s'", out[1].Value)
	}
}

func TestTransform_QuoteAll(t *testing.T) {
	out, _ := Transform(tEntries(), false, TransformQuoteAll)
	if out[0].Value != `"localhost"` {
		t.Errorf("expected quoted value, got %s", out[0].Value)
	}
	// already quoted should not double-quote
	if out[3].Value != `"true"` {
		t.Errorf("expected single-quoted true, got %s", out[3].Value)
	}
}

func TestTransform_UnquoteAll(t *testing.T) {
	out, _ := Transform(tEntries(), false, TransformUnquoteAll)
	if out[3].Value != "true" {
		t.Errorf("expected unquoted true, got %s", out[3].Value)
	}
	// not quoted should remain unchanged
	if out[0].Value != "localhost" {
		t.Errorf("expected localhost unchanged, got %s", out[0].Value)
	}
}

func TestTransform_MultipleOpts(t *testing.T) {
	out, _ := Transform(tEntries(), false, TransformTrimSpace, TransformLowercase)
	if out[1].Value != "5432" {
		t.Errorf("expected 5432 after trim+lower, got %s", out[1].Value)
	}
	if out[2].Value != "myapp" {
		t.Errorf("expected myapp, got %s", out[2].Value)
	}
}

func TestTransform_EmptyEntries(t *testing.T) {
	_, res := Transform([]Entry{}, false, TransformUppercase)
	if res.Total != 0 || res.Modified != 0 {
		t.Errorf("expected zeros for empty input")
	}
}
