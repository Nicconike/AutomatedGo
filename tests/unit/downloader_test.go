package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Nicconike/goautomate/pkg"
	"github.com/schollz/progressbar/v3"
)

func TestDownloadGo(t *testing.T) {
	// Mock the progress bar
	oldNewProgressBar := pkg.NewProgressBar
	pkg.NewProgressBar = func(int64, ...progressbar.Option) *progressbar.ProgressBar {
		return progressbar.NewOptions64(0, progressbar.OptionSetWriter(io.Discard))
	}
	defer func() { pkg.NewProgressBar = oldNewProgressBar }()

	tests := []struct {
		name        string
		version     string
		targetOS    string
		arch        string
		serverResp  string
		statusCode  int
		expectError bool
	}{
		{
			name:        "Successful download",
			version:     "1.16.5",
			targetOS:    "linux",
			arch:        "amd64",
			serverResp:  "mock go binary data",
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "Unsupported OS",
			version:     "1.16.5",
			targetOS:    "unsupported",
			arch:        "amd64",
			serverResp:  "",
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name:        "Unsupported architecture",
			version:     "1.16.5",
			targetOS:    "linux",
			arch:        "unsupported",
			serverResp:  "",
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name:        "Server error",
			version:     "1.16.5",
			targetOS:    "linux",
			arch:        "amd64",
			serverResp:  "Internal Server Error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(50 * time.Millisecond)
				data := make([]byte, 1024)
				w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
				w.WriteHeader(tt.statusCode)
				w.Write(data)
			}))
			defer server.Close()

			// Override the download URL for testing
			originalURL := pkg.DownloadURLFormat
			pkg.DownloadURLFormat = server.URL + "/go%s.%s-%s.%s"
			defer func() { pkg.DownloadURLFormat = originalURL }()

			// Run the download function with a timeout
			errChan := make(chan error, 1)
			go func() {
				errChan <- pkg.DownloadGo(tt.version, tt.targetOS, tt.arch)
			}()

			var err error
			select {
			case err = <-errChan:
			case <-time.After(5 * time.Second):
				t.Fatal("DownloadGo timed out")
			}

			// Check for expected error
			if (err != nil) != tt.expectError {
				t.Errorf("DownloadGo() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// For successful downloads, check if file was created and then remove it
			if !tt.expectError {
				filename := fmt.Sprintf("go%s.%s-%s.tar.gz", tt.version, tt.targetOS, tt.arch)
				if tt.targetOS == "windows" {
					filename = fmt.Sprintf("go%s.%s-%s.zip", tt.version, tt.targetOS, tt.arch)
				}
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					t.Errorf("Expected file %s to be created, but it doesn't exist", filename)
				} else {
					os.Remove(filename)
				}
			}

			// Check for specific error messages
			if tt.expectError {
				switch {
				case tt.targetOS == "unsupported":
					if !strings.Contains(err.Error(), "unsupported operating system") {
						t.Errorf("Expected 'unsupported operating system' error, got: %v", err)
					}
				case tt.arch == "unsupported":
					if !strings.Contains(err.Error(), "unsupported architecture") {
						t.Errorf("Expected 'unsupported architecture' error, got: %v", err)
					}
				case tt.statusCode == http.StatusInternalServerError:
					if !strings.Contains(err.Error(), "unexpected status code: 500") {
						t.Errorf("Expected 'unexpected status code: 500' error, got: %v", err)
					}
				}
			}
		})
	}
}
