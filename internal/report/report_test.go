package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestSummarize(t *testing.T) {
	data := []byte(`{"Results":[{"Target":"alpine:3.19","Vulnerabilities":[{"Severity":"HIGH"},{"Severity":"LOW"},{"Severity":"HIGH"}]}]}`)

	var buf bytes.Buffer
	if err := Summarize(data, &buf); err != nil {
		t.Fatalf("Summarize() error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Total: 3") {
		t.Errorf("expected total 3, got:\n%s", out)
	}
	if !strings.Contains(out, "TOTAL SEVERITY") {
		t.Errorf("expected footer, got:\n%s", out)
	}
	if !strings.Contains(out, "alpine:3.19") {
		t.Errorf("expected target row, got:\n%s", out)
	}
	if !strings.Contains(out, "HIGH") || !strings.Contains(out, "LOW") {
		t.Errorf("expected severity columns in header, got:\n%s", out)
	}
}

func TestSummarizeInvalidJSON(t *testing.T) {
	if err := Summarize([]byte("not json"), &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
