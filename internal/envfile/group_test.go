package envfile

import (
	"testing"
)

func groupEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "STANDALONE", Value: "yes"},
	}
}

func TestGroup_ByPrefix(t *testing.T) {
	result := Group(groupEntries(), GroupOptions{})

	if len(result.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(result.Groups["DB"]))
	}
	if len(result.Groups["APP"]) != 2 {
		t.Errorf("expected 2 APP entries, got %d", len(result.Groups["APP"]))
	}
	if len(result.Groups["other"]) != 1 {
		t.Errorf("expected 1 ungrouped entry, got %d", len(result.Groups["other"]))
	}
}

func TestGroup_OrderedKeys(t *testing.T) {
	result := Group(groupEntries(), GroupOptions{})
	if result.Ordered[0] != "APP" {
		t.Errorf("expected first group APP, got %s", result.Ordered[0])
	}
	if result.Ordered[1] != "DB" {
		t.Errorf("expected second group DB, got %s", result.Ordered[1])
	}
}

func TestGroup_CustomSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "DB-HOST", Value: "localhost"},
		{Key: "DB-PORT", Value: "5432"},
		{Key: "APP-NAME", Value: "myapp"},
	}
	result := Group(entries, GroupOptions{Separator: "-"})
	if len(result.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(result.Groups["DB"]))
	}
}

func TestGroup_CustomUngroupedLabel(t *testing.T) {
	result := Group(groupEntries(), GroupOptions{Ungrouped: "misc"})
	if _, ok := result.Groups["misc"]; !ok {
		t.Error("expected 'misc' group for ungrouped entries")
	}
}

func TestGroup_Depth2(t *testing.T) {
	entries := []Entry{
		{Key: "AWS_S3_BUCKET", Value: "my-bucket"},
		{Key: "AWS_S3_REGION", Value: "us-east-1"},
		{Key: "AWS_EC2_AMI", Value: "ami-123"},
	}
	result := Group(entries, GroupOptions{Depth: 2})
	if len(result.Groups["AWS_S3"]) != 2 {
		t.Errorf("expected 2 AWS_S3 entries, got %d", len(result.Groups["AWS_S3"]))
	}
	if len(result.Groups["AWS_EC2"]) != 1 {
		t.Errorf("expected 1 AWS_EC2 entry, got %d", len(result.Groups["AWS_EC2"]))
	}
}

func TestFormatGroupSummary(t *testing.T) {
	result := Group(groupEntries(), GroupOptions{})
	summary := FormatGroupSummary(result)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	if !containsString(summary, "[DB]") {
		t.Error("expected summary to contain [DB]")
	}
	if !containsString(summary, "[APP]") {
		t.Error("expected summary to contain [APP]")
	}
}
