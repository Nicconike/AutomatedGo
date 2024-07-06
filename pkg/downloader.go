package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

var DownloadURLFormat = "https://dl.google.com/go/go%s.%s-%s.%s"
var validPlatforms = map[string][]string{
	"windows": {"386", "amd64"},
	"linux":   {"386", "amd64", "arm64", "armv6l"},
	"darwin":  {"amd64", "arm64"},
}

func DownloadGo(version, targetOS, arch string) error {
	logger := log.New(os.Stdout, "DownloadGo: ", log.Lshortfile)

	version = strings.TrimPrefix(version, "go")
	logger.Printf("Preparing to download Go version %s", version)

	if targetOS == "" {
		targetOS = runtime.GOOS
		logger.Printf("Target OS not specified, using current OS: %s", targetOS)
	}

	validArchs, ok := validPlatforms[targetOS]
	if !ok {
		return fmt.Errorf("unsupported operating system: %s", targetOS)
	}

	if arch == "" {
		switch targetOS {
		case "windows", "darwin":
			arch = "amd64"
		case "linux":
			arch = runtime.GOARCH
		}
		logger.Printf("Architecture not specified, using default: %s", arch)
	}

	valid := false
	for _, validArch := range validArchs {
		if arch == validArch {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unsupported architecture %s for OS %s", arch, targetOS)
	}

	extension := "tar.gz"
	if targetOS == "windows" {
		extension = "zip"
	}

	url := fmt.Sprintf(DownloadURLFormat, version, targetOS, arch, extension)
	logger.Printf("Download URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	filename := fmt.Sprintf("go%s.%s-%s.%s", version, targetOS, arch, extension)
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Downloading:"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("Download complete!")
		}),
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	logger.Printf("Successfully downloaded Go %s for %s-%s to %s", version, targetOS, arch, filename)
	return nil
}
