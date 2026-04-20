package envfile

import (
	"testing"
)

func lintEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestLint_CleanEntries(t *testing.T) {
	entries := lintEntries("APP_ENV", "production", "DB_HOST", "localhost")
	result := Lint(entries)
	if result.HasIssues() {
		t.Errorf("expected no issues, got: %v", result.Issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	entries := lintEntries("app_env", "production")
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issue for lowercase key")
	}
	if result.Issues[0].Message != "key is not uppercase" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	entries := lintEntries("FOO", "bar", "FOO", "baz")
	result := Lint(entries)
	found := false
	for _, issue := range result.Issues {
		if issue.Message == "duplicate key" {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate key issue")
	}
}

func TestLint_EmptyValue(t *testing.T) {
	entries := lintEntries("MY_VAR", "")
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issue for empty value")
	}
	if result.Issues[0].Message != "value is empty" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_UnquotedSpaces(t *testing.T) {
	entries := lintEntries("MY_VAR", "hello world")
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issue for unquoted value with spaces")
	}
	if result.Issues[0].Message != "value contains spaces but is not quoted" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_InvalidKeyStart(t *testing.T) {
	entries := lintEntries("1BAD", "value")
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issue for key starting with digit")
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{Key: "FOO", Message: "some issue"}
	expected := "[FOO] some issue"
	if issue.String() != expected {
		t.Errorf("expected %q, got %q", expected, issue.String())
	}
}
