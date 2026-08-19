package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"--version"}, strings.NewReader(""), &buf); err != nil {
		t.Fatalf("run(--version) error = %v", err)
	}
	if got, want := buf.String(), "kafka dev\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRunMissingFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "no args",
			args:    nil,
			wantErr: "topic is required",
		},
		{
			name:    "missing topic",
			args:    []string{"--brokers=localhost:9092"},
			wantErr: "topic is required",
		},
		{
			name:    "missing brokers",
			args:    []string{"--topic=trivy"},
			wantErr: "brokers is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args, strings.NewReader(""), &bytes.Buffer{})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("got error %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}
