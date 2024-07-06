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
)

func TestDownloadGo(t *testing.T) {
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
			serverResp:  "pkg Server Error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Logf("Test server received request for URL: %s", r.URL.String())
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResp))
			}))
			defer server.Close()

			t.Logf("Test server started at %s", server.URL)

			// Override the download URL for testing
			originalURL := pkg.DownloadURLFormat
			pkg.DownloadURLFormat = server.URL + "/go%s.%s-%s.%s"
			defer func() { pkg.DownloadURLFormat = originalURL }()

			// Redirect log output
			oldOutput := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Use a channel to signal when the download is complete
			done := make(chan bool)
			var err error

			go func() {
				t.Log("Calling DownloadGo function")
				err = pkg.DownloadGo(tt.version, tt.targetOS, tt.arch)
				t.Log("DownloadGo function returned")
				close(done)
			}()

			// Wait for the download to complete or timeout
			select {
			case <-done:
				t.Log("Download completed")
			case <-time.After(30 * time.Second):
				t.Fatal("Test timed out after 30 seconds")
			}

			// Restore log output
			w.Close()
			os.Stdout = oldOutput
			logOutput, _ := io.ReadAll(r)

			if (err != nil) != tt.expectError {
				t.Errorf("DownloadGo() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				filename := fmt.Sprintf("go%s.%s-%s.tar.gz", tt.version, tt.targetOS, tt.arch)
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					t.Errorf("Expected file %s to be created, but it doesn't exist", filename)
				} else {
					os.Remove(filename) // Clean up the file after test
				}

				if !contains(string(logOutput), "Successfully downloaded") {
					t.Errorf("Expected log output to contain 'Successfully downloaded', but got: %s", logOutput)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
