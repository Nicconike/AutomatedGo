package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func confirmDownload(input io.Reader, output io.Writer) bool {
	reader := bufio.NewReader(input)
	fmt.Fprint(output, "Do you want to download the latest version? (yes/no): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes"
}

func GetDownloadPath(input io.Reader, output io.Writer) string {
	reader := bufio.NewReader(input)
	for {
		fmt.Fprint(output, "Enter the path where you want to download the file (press Enter for current directory, or 'cancel' to abort): ")
		pathChoice, _ := reader.ReadString('\n')
		pathChoice = strings.TrimSpace(pathChoice)

		if pathChoice == "cancel" {
			return ""
		}

		if pathChoice == "" {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(output, "Error getting current directory:", err)
				continue
			}
			fmt.Fprintf(output, "Using current directory: %s\n", currentDir)
			return currentDir
		}

		if _, err := os.Stat(pathChoice); os.IsNotExist(err) {
			fmt.Fprintln(output, "Specified path does not exist. Please try again.")
			continue
		}

		return pathChoice
	}
}

func Run(service VersionChecker, versionFile, currentVersion, targetOS, targetArch string, input io.Reader, output io.Writer) error {
	if versionFile == "" && currentVersion == "" {
		return fmt.Errorf("error: Either -file (-f) or -version (-v) must be specified")
	}

	cv, err := service.GetCurrentVersion(versionFile, currentVersion)
	if err != nil {
		return fmt.Errorf("error getting current version: %v", err)
	}
	fmt.Fprintf(output, "Current version: %s\n", cv)

	latestVersion, err := service.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("error checking latest version: %v", err)
	}
	fmt.Fprintf(output, "Latest version: %s\n", latestVersion)

	if service.IsNewer(latestVersion, cv) {
		fmt.Fprintln(output, "A newer version is available")
		if confirmDownload(input, output) {
			downloadPath := GetDownloadPath(input, output)
			if downloadPath == "" {
				fmt.Fprintln(output, "Download cancelled by user")
				return nil
			}
			err := service.DownloadGo(latestVersion, targetOS, targetArch, downloadPath, input, output)
			if err != nil {
				return fmt.Errorf("error downloading Go: %v", err)
			}
			fmt.Fprintf(output, "%s has been downloaded to %s\n", latestVersion, downloadPath)
		} else {
			fmt.Fprintln(output, "Download aborted by user")
		}
	} else {
		fmt.Fprintln(output, "You have the latest version")
	}
	return nil
}
