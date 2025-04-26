package commands

import (
	"testing"
)

func TestExtractDir(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// the underlaying function retuns . if filepatch is empty as . means current dir
		{"empty", "", "."},
		{"relative", "foo/bar.go", "foo"},
		{"absolute", "/tmp/test/file.go", "/tmp/test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDir(tt.input)
			if got != tt.want {
				t.Errorf("extractDir(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
