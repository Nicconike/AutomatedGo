/*
Package goautomate provides tools for automating Go version checks and downloads.

Features:

  - Detect current Go version from various file types (Dockerfile, go.mod, JSON configs, etc.)
  - Check for the latest available Go version
  - Download the latest Go version if an update is available
  - Support for different operating systems and architectures
  - Checksum validation for downloaded Go versions to ensure integrity

Do's:

 1. Always check for errors returned by functions in this package.
 2. Use GetLatestVersion() to fetch the most recent Go version information.
 3. Specify the correct target OS and architecture when using DownloadGo().
 4. Ensure you have necessary permissions to download and write files.
 5. Define Go version in your project files using one of these formats:
    - In go.mod: go 1.x
    - In Dockerfile:
    * FROM golang:1.x.x
    * ENV GO_VERSION=1.x.x
    - In other files:
    * go_version = "1.x.x"
    * GO_VERSION: 1.x.x
    * golang_version: "1.x.x"
 6. Use the package to automate version checks in your CI/CD pipelines.
 7. Verify checksums of downloaded Go versions for security.

Don'ts:

 1. Don't assume the latest version is always compatible with your project.
 2. Avoid using this package to modify your system's Go installation directly.
 3. Don't use this package in production environments without thorough testing.
 4. Don't ignore version constraints specified in your go.mod file.
 5. Avoid manually modifying files downloaded by this package.
 6. Don't use non-standard formats for specifying Go versions in your project files.

Example usage:

	// Get the latest Go version
	latestVersion, err := pkg.GetLatestVersion()
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Printf("Latest Go version: %s\n", latestVersion)

	// Get current version from a file
	currentVersion, err := pkg.GetCurrentVersion("go.mod", "")
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Printf("Current Go version: %s\n", currentVersion)

	// Check if update is needed
	if pkg.IsNewer(latestVersion, currentVersion) {
	    fmt.Println("An update is available")

	    // Download the latest Go version
	    err = pkg.DownloadGo(latestVersion, "linux", "amd64")
	    if err != nil {
	        log.Fatal(err)
	    }
	    fmt.Println("Successfully downloaded new Go version")

	    // Verify checksum
	    filename := fmt.Sprintf("go%s.linux-amd64.tar.gz", latestVersion)
	    checksum, err := pkg.CalculateFileChecksum(filename)
	    if err != nil {
	        log.Fatal(err)
	    }
	    fmt.Printf("Checksum of downloaded file: %s\n", checksum)
	} else {
	    fmt.Println("You have the latest version")
	}

Functions:

  - GetLatestVersion() (string, error)
    Fetches the latest available Go version.

  - GetCurrentVersion(filename, version string) (string, error)
    Detects the current Go version from a file or uses the provided version.

  - IsNewer(version1, version2 string) bool
    Compares two version strings and returns true if version1 is newer.

  - DownloadGo(version, targetOS, arch string) error
    Downloads the specified Go version for the given OS and architecture.

  - CalculateFileChecksum(filename string) (string, error)
    Calculates the SHA256 checksum of the specified file.

For more detailed information and advanced usage, refer to the README.md file
and the package documentation at https://pkg.go.dev/github.com/Nicconike/goautomate.
*/
package goautomate
