package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nicconike/goautomate/pkg"
)

func TestGetOfficialChecksum(t *testing.T) {
	// Mock server to simulate the Go downloads JSON API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases := []pkg.GoRelease{
			{
				Version: "go1.22.5",
				Files: []struct {
					Filename string `json:"filename"`
					OS       string `json:"os"`
					Arch     string `json:"arch"`
					Version  string `json:"version"`
					SHA256   string `json:"sha256"`
				}{
					{
						Filename: "go1.22.5.linux-amd64.tar.gz",
						SHA256:   "904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	// Replace the actual URL with the test server URL
	originalURL := pkg.URL
	pkg.URL = server.URL
	defer func() { pkg.URL = originalURL }()

	tests := []struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		{"Valid filename", "go1.22.5.linux-amd64.tar.gz", "904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0", false},
		{"Invalid filename", "invalid.tar.gz", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pkg.GetOfficialChecksum(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOfficialChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetOfficialChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateFileChecksum(t *testing.T) {
	// Create a temporary file
	content := []byte("test content")
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

	// Calculate expected SHA256
	expectedSHA256 := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"

	got, err := pkg.CalculateFileChecksum(tmpfile.Name())
	if err != nil {
		t.Fatalf("CalculateFileChecksum() error = %v", err)
	}
	if got != expectedSHA256 {
		t.Errorf("CalculateFileChecksum() = %v, want %v", got, expectedSHA256)
	}

	// Test with non-existent file
	_, err = pkg.CalculateFileChecksum("non_existent_file")
	if err == nil {
		t.Error("CalculateFileChecksum() expected error for non-existent file, got nil")
	}
}
