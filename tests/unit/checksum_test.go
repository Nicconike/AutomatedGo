package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Nicconike/AutomatedGo/v2/pkg"
)

func createServerFunc(filename, sha256 string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
						Filename: filename,
						SHA256:   sha256,
					},
				},
			},
		}
		if err := json.NewEncoder(w).Encode(releases); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

func assertChecksumResult(t *testing.T, got string, err error, want string, wantErr string) {
	if wantErr != "" {
		if err == nil || !strings.Contains(err.Error(), wantErr) {
			t.Errorf("GetOfficialChecksum() error = %v, wantErr %v", err, wantErr)
		}
		return
	}
	if err != nil {
		t.Errorf("GetOfficialChecksum() unexpected error = %v", err)
		return
	}
	if got != want {
		t.Errorf("GetOfficialChecksum() = %v, want %v", got, want)
	}
}

func TestGetOfficialChecksum(t *testing.T) {
	goBinary := "go1.22.5.linux-amd64.tar.gz"

	tests := []struct {
		name       string
		serverFunc func(http.ResponseWriter, *http.Request)
		filename   string
		want       string
		wantErr    string
	}{
		{
			name:       "Valid filename",
			serverFunc: createServerFunc(goBinary, "904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0"),
			filename:   goBinary,
			want:       "904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0",
			wantErr:    "",
		},
		{
			name:       "Invalid filename",
			serverFunc: createServerFunc(goBinary, "904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0"),
			filename:   "invalid.tar.gz",
			want:       "",
			wantErr:    "checksum not found for invalid.tar.gz",
		},
		{
			name: "HTTP failure",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			},
			filename: goBinary,
			want:     "",
			wantErr:  "failed to fetch Go releases: HTTP status 503",
		},
		{
			name: "JSON parsing error",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			filename: goBinary,
			want:     "",
			wantErr:  "failed to parse JSON",
		},
		{
			name: "Read body error",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "{")
				// Simulate a read body error by closing the connection prematurely
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
				conn, _, _ := w.(http.Hijacker).Hijack()
				conn.Close()
			},
			filename: goBinary,
			want:     "",
			wantErr:  "failed to read response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverFunc))
			defer server.Close()

			calculator := pkg.DefaultChecksumCalculator{}
			pkg.URL = server.URL

			got, err := calculator.GetOfficialChecksum(tt.filename)
			assertChecksumResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func createTempFileWithContent(t *testing.T, content string) (*os.File, string) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	expectedSHA256 := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"
	return tmpfile, expectedSHA256
}

func assertChecksum(t *testing.T, got string, err error, expected string) {
	t.Helper()
	if err != nil {
		t.Fatalf("CalculateFileChecksum() error = %v", err)
	}
	if got != expected {
		t.Errorf("CalculateFileChecksum() = %v ,want %v", got, expected)
	}
}

func assertErrorForNonExistentFile(t *testing.T, calculator pkg.DefaultChecksumCalculator) {
	t.Helper()
	_, err := calculator.Calculate("non_existent_file")
	if err == nil {
		t.Error("CalculateFileChecksum() expected error for non-existent file, got nil")
	}
}

func assertErrorForDirectory(t *testing.T, calculator pkg.DefaultChecksumCalculator) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = calculator.Calculate(tmpDir)
	if err == nil {
		t.Error("CalculateFileChecksum() expected error for directory, got nil")
	}
}

func assertErrorForInaccessibleFile(t *testing.T, calculator pkg.DefaultChecksumCalculator) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	os.Remove(tmpfile.Name())

	_, err = calculator.Calculate(tmpfile.Name())
	if err == nil {
		t.Error("CalculateFileChecksum() expected error for non-existent file, got nil")
	}
}

func TestCalculateFileChecksum(t *testing.T) {
	calculator := pkg.DefaultChecksumCalculator{}

	t.Run("Valid file", func(t *testing.T) {
		tmpfile, expectedSHA256 := createTempFileWithContent(t, "test content")
		defer os.Remove(tmpfile.Name())

		got, err := calculator.Calculate(tmpfile.Name())
		assertChecksum(t, got, err, expectedSHA256)
	})

	t.Run("Non-existent file", func(t *testing.T) {
		assertErrorForNonExistentFile(t, calculator)
	})

	t.Run("Directory instead of file", func(t *testing.T) {
		assertErrorForDirectory(t, calculator)
	})

	t.Run("File becomes inaccessible", func(t *testing.T) {
		assertErrorForInaccessibleFile(t, calculator)
	})
}
