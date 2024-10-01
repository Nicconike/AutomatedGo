package pkg

import (
	"errors"
	"fmt"
	"io"
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

type DefaultDownloader struct{}

func (d *DefaultDownloader) Download(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
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
			fmt.Printf("\nDownload Complete!")
		}),
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %w", err)
	}

	return nil
}

type DefaultRemover struct{}

func (r *DefaultRemover) Remove(filename string) error {
	return os.Remove(filename)
}

func DownloadGo(version, targetOS, arch, path string, downloader FileDownloader, remover FileRemover, checksum ChecksumCalculator) error {
	version = strings.TrimPrefix(version, "go")
	fmt.Printf("Preparing to download Go version %s\n", version)

	if targetOS == "" {
		targetOS = runtime.GOOS
		fmt.Printf("Target OS not specified, using current OS: %s\n", targetOS)
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
		default:
			return fmt.Errorf("unsupported operating system: %s", targetOS)
		}
		fmt.Printf("Architecture not specified, using default: %s\n", arch)
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

	filename := fmt.Sprintf("go%s.%s-%s.%s", version, targetOS, arch, extension)
	fmt.Printf("Fetching Official Checksum for %s\n", filename)

	officialChecksum, err := checksum.GetOfficialChecksum(filename)
	if err != nil {
		fmt.Printf("Failed to get official checksum: %s\n", err)
		return err
	} else {
		fmt.Printf("Successfully fetched official checksum: %s\n", officialChecksum)
	}

	url := fmt.Sprintf(DownloadURLFormat, version, targetOS, arch, extension)
	err = downloader.Download(url, filename)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}

	fmt.Printf("\nCalculating checksum for %s\n", filename)
	calculatedChecksum, err := checksum.Calculate(filename)
	if err != nil || calculatedChecksum != officialChecksum {
		if removeErr := remover.Remove(filename); removeErr != nil {
			fmt.Printf("Error removing file %s after failed checksum calculation: %s\n", filename, removeErr)
		}
		errMsg := fmt.Sprintf("Checksum mismatch: expected %s, got %s for file %s", officialChecksum, calculatedChecksum, filename)
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	fmt.Println("Checksum verification successful!")
	return nil
}
