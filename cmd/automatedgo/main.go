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
	targetArch := flag.String("arch", "", "Target architecture (386, amd64, armv6l)")

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

	// Initialize the VersionService with default implementations
	service := &pkg.VersionService{
		Downloader: &pkg.DefaultDownloader{},
		Remover:    &pkg.DefaultRemover{},
		Checksum:   &pkg.DefaultChecksumCalculator{},
	}

	if err := pkg.Run(service, *versionFile, *currentVersion, *targetOS, *targetArch, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
