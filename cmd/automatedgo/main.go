package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Nicconike/AutomatedGo/pkg"
)

func main() {
	// Define flags
	versionFile := flag.String("file", "", "Path to file containing current Go version")
	currentVersion := flag.String("version", "", "Current Go version")
	targetOS := flag.String("os", "", "Target operating system (windows, linux, darwin)")
	targetArch := flag.String("arch", "", "Target architecture (386, amd64, arm64, armv6l)")

	// Add aliases for short versions
	flag.StringVar(versionFile, "f", "", "Path to file containing current Go version (shorthand)")
	flag.StringVar(currentVersion, "v", "", "Current Go version (shorthand)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [-os=<OS>] [-arch=<ARCH>] (-file|-f=<path> | -version|-v=<version>)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check if either file or version is specified
	if *versionFile == "" && *currentVersion == "" {
		fmt.Println("Error: Either -file (-f) or -version (-v) must be specified")
		flag.Usage()
		os.Exit(1)
	}

	cv, err := pkg.GetCurrentVersion(*versionFile, *currentVersion)
	if err != nil {
		fmt.Printf("Error getting current version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current version: %s\n", cv)

	latestVersion, err := pkg.GetLatestVersion()
	if err != nil {
		fmt.Printf("Error checking latest version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Latest version: %s\n", latestVersion)

	if pkg.IsNewer(latestVersion, cv) {
		fmt.Println("A newer version is available")
		err := pkg.DownloadGo(latestVersion, *targetOS, *targetArch)
		if err != nil {
			fmt.Printf("Error downloading Go: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("You have the latest version")
	}
}
