package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	versionURL        = "https://go.dev/VERSION?m=text"
	downloadURLFormat = "https://dl.google.com/go/go%s.%s-%s.tar.gz"
	checksumURLFormat = "https://dl.google.com/go/go%s.%s-%s.tar.gz.sha256"
)

func main() {
	// Check latest version
	latestVersion := getLatestVersion()
	currentVersion := getCurrentVersion()

	if latestVersion == currentVersion {
		fmt.Println("Go is already up to date")
		return
	}

	fmt.Printf("Updating Go from %s to %s\n", currentVersion, latestVersion)

	// Construct download URLs
	osType := runtime.GOOS
	arch := runtime.GOARCH
	tarURL := fmt.Sprintf(downloadURLFormat, latestVersion, osType, arch)
	checksumURL := fmt.Sprintf(checksumURLFormat, latestVersion, osType, arch)
	tarFilename := fmt.Sprintf("go%s.%s-%s.tar.gz", latestVersion, osType, arch)

	// Download Go binary
	downloadFile(tarURL, tarFilename)

	// Verify checksum
	expectedChecksum := downloadChecksum(checksumURL)
	actualChecksum := calculateChecksum(tarFilename)

	if expectedChecksum != actualChecksum {
		fmt.Println("Checksum verification failed")
		os.Remove(tarFilename)
		return
	}

	fmt.Println("Checksum verification passed")

	// Upload to Git LFS
	uploadToGitLFS(tarFilename)

	// Clean up
	os.Remove(tarFilename)

	fmt.Println("Go binary downloaded and uploaded to Git LFS successfully")
}

func getLatestVersion() string {
	resp, err := http.Get(versionURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(body))
}

func getCurrentVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	parts := strings.Split(string(output), " ")
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "go")
	}
	return "unknown"
}

func downloadFile(url string, filepath string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Downloaded %s\n", filepath)
}

func downloadChecksum(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(body), " ")[0]
}

func calculateChecksum(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func uploadToGitLFS(filepath string) {
	// Ensure Git LFS is installed and initialized
	exec.Command("git", "lfs", "install").Run()

	// Track .tar.gz files with Git LFS
	trackCmd := exec.Command("git", "lfs", "track", "*.tar.gz")
	trackCmd.Run()

	// Add the file to Git
	addCmd := exec.Command("git", "add", filepath)
	addCmd.Run()

	// Commit the file
	commitCmd := exec.Command("git", "commit", "-m", fmt.Sprintf("Add Go binary %s", filepath))
	commitCmd.Run()

	// Push the changes
	pushCmd := exec.Command("git", "push")
	pushCmd.Run()

	fmt.Printf("Uploaded %s to Git LFS\n", filepath)
}
