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

func TestShouldProcessFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		testType string
		want     bool
	}{
		// Golang files
		{"golang file with .go extension", "main.go", "golang", true},
		{"golang file with path", "cmd/rgt/main.go", "golang", true},
		{"golang file with absolute path", "/home/user/project/main.go", "golang", true},
		{"non-go file with golang type", "README.md", "golang", false},
		{"python file with golang type", "test.py", "golang", false},
		{"txt file with golang type", "file.txt", "golang", false},

		// Python files
		{"python file with .py extension", "test.py", "python", true},
		{"python file with path", "tests/test_main.py", "python", true},
		{"python file with absolute path", "/home/user/project/test.py", "python", true},
		{"non-python file with python type", "README.md", "python", false},
		{"go file with python type", "main.go", "python", false},
		{"txt file with python type", "file.txt", "python", false},

		// Unknown test type (backward compatible)
		{"go file with unknown type", "main.go", "rust", true},
		{"py file with unknown type", "test.py", "rust", true},
		{"any file with unknown type", "README.md", "rust", true},

		// Edge cases
		{"file without extension golang", "Makefile", "golang", false},
		{"file without extension python", "Dockerfile", "python", false},
		{"hidden go file", ".hidden.go", "golang", true},
		{"hidden py file", ".hidden.py", "python", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldProcessFile(tt.filePath, tt.testType)
			if got != tt.want {
				t.Errorf("shouldProcessFile(%q, %q) = %v, want %v", tt.filePath, tt.testType, got, tt.want)
			}
		})
	}
}
