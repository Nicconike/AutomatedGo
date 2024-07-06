package tests

import (
	"os"
	"testing"

	"github.com/Nicconike/goautomate/pkg"
)

func TestGetCurrentVersion(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		directVersion  string
		expectedResult string
		expectError    bool
	}{
		{"Direct version", "", "1.16.5", "1.16.5", false},
		{"No input", "", "", "", true},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkg.GetCurrentVersion(tt.filePath, tt.directVersion)
			if (err != nil) != tt.expectError {
				t.Errorf("GetCurrentVersion() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if result != tt.expectedResult {
				t.Errorf("GetCurrentVersion() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestExtractGoVersion(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{"Dockerfile", "FROM golang:1.16.5", "1.16.5"},
		{"go.mod", "go 1.17", "1.17"},
		{"JSON", `{"go_version": "1.18.0"}`, "1.18.0"},
		{"No version", "Some random content", ""},
		// Add more test cases for different formats
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.ExtractGoVersion(tt.content)
			if result != tt.expected {
				t.Errorf("ExtractGoVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadVersionFromFile(t *testing.T) {
	// Create a temporary file for testing
	content := []byte("go 1.17")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading from the file
	version, err := pkg.ReadVersionFromFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if version != "1.17" {
		t.Errorf("Expected version 1.17, got %s", version)
	}

	// Test reading from non-existent file
	_, err = pkg.ReadVersionFromFile("non_existent_file.txt")
	if err == nil {
		t.Error("Expected an error for non-existent file, got nil")
	}
}
