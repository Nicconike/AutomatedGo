package pkg

import (
	"io"
	"net/http"
	"strings"
)

var VersionURL = "https://go.dev/VERSION?m=text"

func GetLatestVersion() (string, error) {
	resp, err := http.Get(VersionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(body))
	version = strings.Split(version, "\n")[0]
	return version, nil
}
