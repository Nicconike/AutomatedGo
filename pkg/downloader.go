package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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

type DownloadConfig struct {
	Version    string
	TargetOS   string
	Arch       string
	Path       string
	Downloader FileDownloader
	Remover    FileRemover
	Checksum   ChecksumCalculator
	Input      io.Reader
	Output     io.Writer
}

func DownloadGo(config DownloadConfig) error {
	version := strings.TrimPrefix(config.Version, "go")
	fmt.Fprintf(config.Output, "Preparing to download Go version %s\n", version)

	if config.TargetOS == "" {
		fmt.Fprint(config.Output, "Enter target OS (windows, linux, darwin): ")
		fmt.Fscan(config.Input, &config.TargetOS)
	}

	validArchs, ok := validPlatforms[config.TargetOS]
	if !ok {
		return fmt.Errorf("unsupported operating system: %s", config.TargetOS)
	}

	if config.Arch == "" {
		fmt.Fprintf(config.Output, "Enter target architecture %v: ", validArchs)
		fmt.Fscan(config.Input, &config.Arch)
	}

	valid := false
	for _, validArch := range validArchs {
		if config.Arch == validArch {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unsupported architecture %s for OS %s", config.Arch, config.TargetOS)
	}

	extension := "tar.gz"
	if config.TargetOS == "windows" {
		extension = "zip"
	}

	filename := fmt.Sprintf("go%s.%s-%s.%s", version, config.TargetOS, config.Arch, extension)
	fmt.Fprintf(config.Output, "Fetching Official Checksum for %s\n", filename)

	officialChecksum, err := config.Checksum.GetOfficialChecksum(filename)
	if err != nil {
		fmt.Fprintf(config.Output, "Failed to get official checksum: %s\n", err)
		return err
	}
	fmt.Fprintf(config.Output, "Successfully fetched official checksum: %s\n", officialChecksum)

	url := fmt.Sprintf(DownloadURLFormat, version, config.TargetOS, config.Arch, extension)
	err = config.Downloader.Download(url, filename)
	if err != nil {
		fmt.Fprintf(config.Output, "Error downloading file: %s\n", err)
		return err
	}

	fmt.Fprintf(config.Output, "\nCalculating checksum for %s\n", filename)
	calculatedChecksum, err := config.Checksum.Calculate(filename)
	if err != nil || calculatedChecksum != officialChecksum {
		if removeErr := config.Remover.Remove(filename); removeErr != nil {
			fmt.Fprintf(config.Output, "Error removing file %s after failed checksum calculation: %s\n", filename, removeErr)
		}
		errMsg := fmt.Sprintf("Checksum mismatch: expected %s, got %s for file %s", officialChecksum, calculatedChecksum, filename)
		fmt.Fprintln(config.Output, errMsg)
		return errors.New(errMsg)
	}

	fmt.Fprintln(config.Output, "Checksum verification successful!")
	return nil
}
