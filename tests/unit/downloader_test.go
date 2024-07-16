package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Nicconike/goautomate/pkg"
	"github.com/schollz/progressbar/v3"
)

// Mock functions
func mockGetOfficialChecksum(filename string) (string, error) {
	return "mockedchecksum", nil
}

func mockCalculateFileChecksum(filename string) (string, error) {
	return "mockedchecksum", nil
}

func TestDownloadFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	// Create a temporary file to download to
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = pkg.DownloadFile(server.URL, tmpfile.Name())
	if err != nil {
		t.Errorf("downloadFile returned an error: %v", err)
	}

	// Check the contents of the downloaded file
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "test content" {
		t.Errorf("downloaded content does not match expected. Got %s, want %s", string(content), "test content")
	}
}

func TestDownloadFileErrors(t *testing.T) {
	// Test server error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := pkg.DownloadFile(server.URL, "testfile")
	if err == nil {
		t.Error("Expected an error for server error, got nil")
	}

	// Test invalid file path
	err = pkg.DownloadFile("http://example.com", "/invalid/path/testfile")
	if err == nil {
		t.Error("Expected an error for invalid file path, got nil")
	}
}
func TestDownloadGo(t *testing.T) {
	// Mock the progress bar
	oldNewProgressBar := pkg.NewProgressBar
	pkg.NewProgressBar = func(max int64, _ ...progressbar.Option) *progressbar.ProgressBar {
		return progressbar.NewOptions64(max, progressbar.OptionSetWriter(io.Discard))
	}
	defer func() { pkg.NewProgressBar = oldNewProgressBar }()

	// Mock GetOfficialChecksum and CalculateFileChecksum
	oldGetOfficialChecksum := pkg.GetOfficialChecksum
	oldCalculateFileChecksum := pkg.CalculateFileChecksum
	pkg.GetOfficialChecksum = mockGetOfficialChecksum
	pkg.CalculateFileChecksum = mockCalculateFileChecksum
	defer func() {
		pkg.GetOfficialChecksum = oldGetOfficialChecksum
		pkg.CalculateFileChecksum = oldCalculateFileChecksum
	}()

	// Mock DownloadFile
	oldDownloadFile := pkg.DownloadFile
	pkg.DownloadFile = func(url, filename string) error {
		if strings.Contains(url, "500") {
			return fmt.Errorf("error downloading: unexpected status code: 500")
		}
		return nil
	}
	defer func() { pkg.DownloadFile = oldDownloadFile }()

	tests := []struct {
		name        string
		version     string
		targetOS    string
		arch        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Successful download",
			version:     "1.22.5",
			targetOS:    "linux",
			arch:        "amd64",
			expectError: false,
		},
		{
			name:        "Unsupported OS",
			version:     "1.22.5",
			targetOS:    "unsupported",
			arch:        "amd64",
			expectError: true,
			errorMsg:    "unsupported operating system: unsupported",
		},
		{
			name:        "Unsupported architecture",
			version:     "1.22.5",
			targetOS:    "linux",
			arch:        "unsupported",
			expectError: true,
			errorMsg:    "unsupported architecture unsupported for OS linux",
		},
		{
			name:        "Server error",
			version:     "1.22.5",
			targetOS:    "linux",
			arch:        "amd64",
			expectError: true,
			errorMsg:    "error downloading: unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the download URL for testing
			originalURL := pkg.DownloadURLFormat
			if tt.name == "Server error" {
				pkg.DownloadURLFormat = "https://example.com/500/go%s.%s-%s.%s"
			} else {
				pkg.DownloadURLFormat = "https://example.com/go%s.%s-%s.%s"
			}
			defer func() { pkg.DownloadURLFormat = originalURL }()

			// Run the download function
			err := pkg.DownloadGo(tt.version, tt.targetOS, tt.arch)

			// Check for expected error
			if (err != nil) != tt.expectError {
				t.Errorf("DownloadGo() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Check for specific error messages
			if tt.expectError {
				if err == nil || !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errorMsg, err)
				}
			}
		})
	}
}
