package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nicconike/AutomatedGo/v2/pkg"
)

func createTestServer(status int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(response))
	}))
}

func createTemp(t *testing.T) *os.File {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "downloaded")
	if err != nil {
		t.Fatal(err)
	}
	return tmpfile
}

func runDownloadTest(t *testing.T, url string, filename string, expectedError string, expectedContent string) {
	t.Helper()
	downloader := pkg.DefaultDownloader{}
	err := downloader.Download(url, filename)

	if expectedError != "" {
		if err == nil || err.Error() != expectedError {
			t.Errorf("Download() error = %v, wantErr %v", err, expectedError)
		}
		return
	}

	if err != nil {
		t.Fatalf("Download() unexpected error = %v", err)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(content) != expectedContent {
		t.Errorf("Downloaded content = %v, want %v", string(content), expectedContent)
	}
}

func TestDownload(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name:           "Successful download",
			serverResponse: "file content",
			serverStatus:   http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Server error",
			serverResponse: "",
			serverStatus:   http.StatusInternalServerError,
			expectedError:  "unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createTestServer(tt.serverStatus, tt.serverResponse)
			defer server.Close()

			tmpfile := createTemp(t)
			defer os.Remove(tmpfile.Name())

			runDownloadTest(t, server.URL, tmpfile.Name(), tt.expectedError, tt.serverResponse)
		})
	}
}
