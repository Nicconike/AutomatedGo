/*
**AutomatedGo Package**

AutomatedGo provides tools for automated Go version management and updates.

Features:

  - Get the current Go version from a specified file or input
  - Check for the latest available Go version
  - Compare Go versions to determine if an update is available
  - Download the latest Go version for different operating systems and architectures
  - Checksum validation for downloaded Go versions to ensure integrity

Usage Example:

package main

import (

	"fmt"
	"log"

	"github.com/Nicconike/AutomatedGo/v2/pkg"

)

	func main() {
	    // Create a new VersionService
	    service := pkg.NewVersionService(
	        &pkg.DefaultDownloader{},
	        &pkg.DefaultRemover{},
	        &pkg.DefaultChecksumCalculator{},
	    )

	    // Get the current Go version
	    currentVersion, err := service.GetCurrentVersion("", "")
	    if err != nil {
	        log.Fatalf("Error getting current version: %v", err)
	    }
	    fmt.Printf("Current Go version: %s\n", currentVersion)

	    // Check for the latest Go version
	    latestVersion, err := service.GetLatestVersion()
	    if err != nil {
	        log.Fatalf("Error getting latest version: %v", err)
	    }
	    fmt.Printf("Latest Go version: %s\n", latestVersion)

	    // Check if an update is available
	    if service.IsNewer(latestVersion, currentVersion) {
	        fmt.Println("An update is available!")

	        // Download the latest version
	        err = service.DownloadGo(latestVersion, "", "", "/tmp", nil, nil)
	        if err != nil {
	            log.Fatalf("Error downloading Go: %v", err)
	        }
	        fmt.Printf("Successfully downloaded Go %s to /tmp\n", latestVersion)
	    } else {
	        fmt.Println("You have the latest version of Go.")
	    }
	}

This example demonstrates how to use the AutomatedGo package to check for updates,
compare versions, and download the latest version of Go if an update is available.

Do's:

 1. Always check for errors returned by functions in this package.
 2. Use GetLatestVersion() to fetch the most recent Go version information.
 3. Use GetCurrentVersion() with appropriate parameters to determine the current version.
 4. Specify the correct target OS and architecture when using DownloadGo() if needed.
 5. Ensure you have necessary permissions to download and write files.
 6. Use IsNewer() to compare versions and determine if an update is needed.
 7. Use the package to automate version checks in your CI/CD pipelines.

Don'ts:

 1. Don't assume the latest version is always compatible with your project.
 2. Avoid using this package to modify your system's Go installation directly.
 3. Don't use this package in production environments without thorough testing.
 4. Don't ignore version constraints specified in your go.mod file.
 5. Avoid manually modifying files downloaded by this package.

Note: The package allows flexible version checking and downloading. You can provide
specific version information or let the package determine versions automatically.
*/
package AutomatedGo
