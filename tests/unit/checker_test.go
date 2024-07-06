package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nicconike/goautomate/pkg"
)

var originalURL string

func setTestVersionURL(url string) {
	originalURL = pkg.VersionURL
	pkg.VersionURL = url
}

func resetTestVersionURL() {
	pkg.VersionURL = originalURL
}

func TestGetLatestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("go1.17.1\n"))
	}))
	defer server.Close()

	setTestVersionURL(server.URL)
	defer resetTestVersionURL()

	version, err := pkg.GetLatestVersion()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if version != "go1.17.1" {
		t.Errorf("Expected version go1.17.1, got %s", version)
	}
}
