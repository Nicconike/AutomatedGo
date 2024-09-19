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

func getDownloadPath(input io.Reader, output io.Writer) string {
	reader := bufio.NewReader(input)
	fmt.Fprint(output, "Do you want to download it to the current working directory? (yes/no): ")
	locationChoice, _ := reader.ReadString('\n')
	locationChoice = strings.TrimSpace(strings.ToLower(locationChoice))

	if locationChoice == "no" {
		fmt.Fprint(output, "Enter the path where you want to download the file: ")
		pathChoice, _ := reader.ReadString('\n')
		pathChoice = strings.TrimSpace(pathChoice)

		if _, err := os.Stat(pathChoice); os.IsNotExist(err) {
			return ""
		}

		return pathChoice
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return currentDir
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
			downloadPath := getDownloadPath(input, output)
			if downloadPath == "" {
				return fmt.Errorf("specified path does not exist")
			}
			err := service.DownloadGo(latestVersion, targetOS, targetArch, downloadPath)
			if err != nil {
				return fmt.Errorf("error downloading Go: %v", err)
			}
		} else {
			fmt.Fprintln(output, "Download aborted by user")
		}
	} else {
		fmt.Fprintln(output, "You have the latest version")
	}
	return nil
}
