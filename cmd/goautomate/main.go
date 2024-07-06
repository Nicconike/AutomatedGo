package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Nicconike/goautomate/pkg"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// Define flags
	versionFile := flag.String("file", "", "Path to file containing current Go version")
	currentVersion := flag.String("version", "", "Current Go version")
	versionShort := flag.String("v", "", "Current Go version (short flag)")
	targetOS := flag.String("os", "", "Target operating system (windows, linux, darwin)")
	targetArch := flag.String("arch", "", "Target architecture (386, amd64, arm64, armv6l)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [-os=<OS>] [-arch=<ARCH>] [-file=<path> | -version=<version> | -v=<version>]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check if either file or version is specified
	if *versionFile == "" && *currentVersion == "" && *versionShort == "" {
		fmt.Fprintln(os.Stderr, "Error: Either -file, -version, or -v must be specified")
		flag.Usage()
		os.Exit(1)
	}

	// Prioritize -version over -v if both are provided
	if *currentVersion == "" {
		*currentVersion = *versionShort
	}

	cv, err := pkg.GetCurrentVersion(*versionFile, *currentVersion)
	if err != nil {
		log.Fatalf("Error getting current version: %v", err)
	}

	log.Printf("Current version: %s", cv)

	latestVersion, err := pkg.GetLatestVersion()
	if err != nil {
		log.Fatalf("Error checking latest version: %v", err)
	}

	log.Printf("Latest version: %s", latestVersion)

	if pkg.IsNewer(latestVersion, cv) {
		log.Println("A newer version is available. Downloading...")
		err := pkg.DownloadGo(latestVersion, *targetOS, *targetArch)
		if err != nil {
			log.Fatalf("Error downloading Go: %v", err)
		}
	} else {
		log.Println("You have the latest version.")
	}
}
