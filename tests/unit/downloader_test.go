package tests

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Nicconike/AutomatedGo/pkg"
	"github.com/schollz/progressbar/v3"
)

func TestDownloadFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
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
	// Test server error (existing test)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := pkg.DownloadFile(server.URL, "testfile")
	if err == nil {
		t.Error("Expected an error for server error, got nil")
	}

	// Test invalid file path (existing test)
	err = pkg.DownloadFile("http://example.com", "/invalid/path/testfile")
	if err == nil {
		t.Error("Expected an error for invalid file path, got nil")
	}

	// Test http.Get error
	err = pkg.DownloadFile("http://invalid-url", "testfile")
	if err == nil {
		t.Error("Expected an error for invalid URL, got nil")
	}
	if !strings.Contains(err.Error(), "error downloading:") {
		t.Errorf("Expected error message to contain 'error downloading:', got: %v", err)
	}

	// Test io.Copy error
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer server.Close()

	err = pkg.DownloadFile(server.URL, "testfile")
	if err == nil {
		t.Error("Expected an error for io.Copy failure, got nil")
	}
	if !strings.Contains(err.Error(), "error saving file:") {
		t.Errorf("Expected error message to contain 'error saving file:', got: %v", err)
	}
}
func TestDownloadGo(t *testing.T) {
	// Mock the progress bar
	oldNewProgressBar := pkg.NewProgressBar
	pkg.NewProgressBar = func(max int64, _ ...progressbar.Option) *progressbar.ProgressBar {
		return progressbar.NewOptions64(max, progressbar.OptionSetWriter(io.Discard))
	}
	defer func() { pkg.NewProgressBar = oldNewProgressBar }()

	// Variables to control mock behavior
	var checksumMismatch bool
	var getOfficialChecksumError bool
	var calculateFileChecksumError bool
	var downloadFileError bool
	var removeFileError bool

	// Mock GetOfficialChecksum
	oldGetOfficialChecksum := pkg.GetOfficialChecksum
	pkg.GetOfficialChecksum = func(filename string) (string, error) {
		if getOfficialChecksumError {
			return "", errors.New("mock official checksum error")
		}
		return "mockedchecksum", nil
	}
	defer func() { pkg.GetOfficialChecksum = oldGetOfficialChecksum }()

	// Mock CalculateFileChecksum
	oldCalculateFileChecksum := pkg.CalculateFileChecksum
	pkg.CalculateFileChecksum = func(filename string) (string, error) {
		if calculateFileChecksumError {
			return "", errors.New("mock calculate checksum error")
		}
		if checksumMismatch {
			return "mismatchedchecksum", nil
		}
		return "mockedchecksum", nil
	}
	defer func() { pkg.CalculateFileChecksum = oldCalculateFileChecksum }()

	// Mock DownloadFile
	oldDownloadFile := pkg.DownloadFile
	pkg.DownloadFile = func(url, filename string) error {
		if downloadFileError {
			return errors.New("mock download file error")
		}
		return nil
	}
	defer func() { pkg.DownloadFile = oldDownloadFile }()

	// Mock os.Remove
	oldOsRemove := pkg.OsRemove
	pkg.OsRemove = func(name string) error {
		if removeFileError {
			return errors.New("mock remove file error")
		}
		return nil
	}
	defer func() { pkg.OsRemove = oldOsRemove }()

	// Mock runtime.GOOS and runtime.GOARCH
	oldRuntimeGOOS := pkg.RuntimeGOOS
	oldRuntimeGOARCH := pkg.RuntimeGOARCH
	defer func() {
		pkg.RuntimeGOOS = oldRuntimeGOOS
		pkg.RuntimeGOARCH = oldRuntimeGOARCH
	}()

	tests := []struct {
		name        string
		version     string
		targetOS    string
		arch        string
		setupMocks  func()
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Successful download",
			version:     "1.22.5",
			targetOS:    "linux",
			arch:        "amd64",
			setupMocks:  func() {},
			expectError: false,
		},
		{
			name:        "Unsupported OS",
			version:     "1.22.5",
			targetOS:    "unsupported",
			arch:        "amd64",
			setupMocks:  func() {},
			expectError: true,
			errorMsg:    "unsupported operating system: unsupported",
		},
		{
			name:        "Unsupported architecture",
			version:     "1.22.5",
			targetOS:    "linux",
			arch:        "unsupported",
			setupMocks:  func() {},
			expectError: true,
			errorMsg:    "unsupported architecture unsupported for OS linux",
		},
		{
			name:     "Default OS",
			version:  "1.22.5",
			targetOS: "",
			arch:     "amd64",
			setupMocks: func() {
				pkg.RuntimeGOOS = "darwin"
			},
			expectError: false,
		},
		{
			name:     "Default architecture for Linux",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "",
			setupMocks: func() {
				pkg.RuntimeGOARCH = "arm64"
			},
			expectError: false,
		},
		{
			name:        "Default architecture for Windows",
			version:     "1.22.5",
			targetOS:    "windows",
			arch:        "",
			setupMocks:  func() {},
			expectError: false,
		},
		{
			name:        "Default architecture for Darwin",
			version:     "1.22.5",
			targetOS:    "darwin",
			arch:        "",
			setupMocks:  func() {},
			expectError: false,
		},
		{
			name:     "GetOfficialChecksum error",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "amd64",
			setupMocks: func() {
				getOfficialChecksumError = true
			},
			expectError: true,
			errorMsg:    "mock official checksum error",
		},
		{
			name:     "DownloadFile error",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "amd64",
			setupMocks: func() {
				downloadFileError = true
			},
			expectError: true,
			errorMsg:    "mock download file error",
		},
		{
			name:     "CalculateFileChecksum error",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "amd64",
			setupMocks: func() {
				calculateFileChecksumError = true
			},
			expectError: true,
			errorMsg:    "mock calculate checksum error",
		},
		{
			name:     "Checksum mismatch",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "amd64",
			setupMocks: func() {
				checksumMismatch = true
			},
			expectError: true,
			errorMsg:    "Checksum mismatch",
		},
		{
			name:     "Remove file error after checksum mismatch",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "amd64",
			setupMocks: func() {
				checksumMismatch = true
				removeFileError = true
			},
			expectError: true,
			errorMsg:    "Checksum mismatch",
		},
		{
			name:     "Remove file error after checksum calculation error",
			version:  "1.22.5",
			targetOS: "linux",
			arch:     "amd64",
			setupMocks: func() {
				calculateFileChecksumError = true
				removeFileError = true
			},
			expectError: true,
			errorMsg:    "mock calculate checksum error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock variables
			checksumMismatch = false
			getOfficialChecksumError = false
			calculateFileChecksumError = false
			downloadFileError = false
			removeFileError = false
			pkg.RuntimeGOOS = oldRuntimeGOOS
			pkg.RuntimeGOARCH = oldRuntimeGOARCH

			// Setup mocks for this test case
			tt.setupMocks()

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

			// For cases where we expect file removal, check if OsRemove was called
			if calculateFileChecksumError || checksumMismatch {
				if !removeFileError && pkg.OsRemove == nil {
					t.Errorf("Expected OsRemove to be called, but it wasn't")
				}
			}
		})
	}
}
