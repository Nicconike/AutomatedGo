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
	tests := []struct {
		name       string
		serverFunc func(http.ResponseWriter, *http.Request)
		filename   string
		want       string
		wantErr    bool
	}{
		{
			name: "Valid filename",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
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
			},
			filename: "go1.22.5.linux-amd64.tar.gz",
			want:     "904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0",
			wantErr:  false,
		},
		{
			name: "Invalid filename",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
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
			},
			filename: "invalid.tar.gz",
			want:     "",
			wantErr:  true,
		},
		{
			name: "HTTP error",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			},
			filename: "go1.22.5.linux-amd64.tar.gz",
			want:     "",
			wantErr:  true,
		},
		{
			name: "Invalid JSON",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("invalid json"))
			},
			filename: "go1.22.5.linux-amd64.tar.gz",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverFunc))
			defer server.Close()

			originalURL := pkg.URL
			pkg.URL = server.URL
			defer func() { pkg.URL = originalURL }()

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
	// Test with a valid file
	t.Run("Valid file", func(t *testing.T) {
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

		expectedSHA256 := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"

		got, err := pkg.CalculateFileChecksum(tmpfile.Name())
		if err != nil {
			t.Fatalf("CalculateFileChecksum() error = %v", err)
		}
		if got != expectedSHA256 {
			t.Errorf("CalculateFileChecksum() = %v, want %v", got, expectedSHA256)
		}
	})

	// Test with a non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		_, err := pkg.CalculateFileChecksum("non_existent_file")
		if err == nil {
			t.Error("CalculateFileChecksum() expected error for non-existent file, got nil")
		}
	})

	// Test with a directory instead of a file
	t.Run("Directory instead of file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "testdir")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		_, err = pkg.CalculateFileChecksum(tmpDir)
		if err == nil {
			t.Error("CalculateFileChecksum() expected error for directory, got nil")
		}
	})

	// Test with a file that becomes inaccessible after opening
	t.Run("File becomes inaccessible", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatal(err)
		}
		tmpfile.Close()
		// Remove the file instead of just changing permissions
		os.Remove(tmpfile.Name())

		_, err = pkg.CalculateFileChecksum(tmpfile.Name())
		if err == nil {
			t.Error("CalculateFileChecksum() expected error for non-existent file, got nil")
		}
	})
}