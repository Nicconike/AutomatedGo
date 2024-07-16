package pkg

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GoRelease struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		OS       string `json:"os"`
		Arch     string `json:"arch"`
		Version  string `json:"version"`
		SHA256   string `json:"sha256"`
	} `json:"files"`
}

var URL = "https://go.dev/dl/?mode=json"

var GetOfficialChecksum = getOfficialChecksum
var CalculateFileChecksum = calculateFileChecksum

func getOfficialChecksum(filename string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Go releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch Go releases: HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var releases []GoRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, release := range releases {
		for _, file := range release.Files {
			if file.Filename == filename {
				return file.SHA256, nil
			}
		}
	}

	return "", fmt.Errorf("checksum not found for %s", filename)
}

func calculateFileChecksum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
