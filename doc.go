/*
Package goautomate provides tools for automating Go version checks and downloads.

Do's:
1. Always check for errors returned by functions in this package.
2. Use GetLatestVersion() to fetch the most recent Go version information.
3. Specify the correct target OS and architecture when using DownloadGo().
4. Ensure you have necessary permissions to download and write files.
5. Define Go version in your project files using one of these formats:
  - In go.mod: go 1.x
  - In Dockerfile:
  - - FROM golang:1.x.x
  - - ENV GO_VERSION=1.x.x
  - In other files:
  - - go_version = "1.x.x"
  - - GO_VERSION: 1.x.x
  - - golang_version: "1.x.x"

6. Use the package to automate version checks in your CI/CD pipelines.

Don'ts:
1. Don't assume the latest version is always compatible with your project.
2. Avoid using this package to modify your system's Go installation directly.
3. Don't use this package in production environments without thorough testing.
4. Don't ignore version constraints specified in your go.mod file.
5. Avoid manually modifying files downloaded by this package.
6. Don't use non-standard formats for specifying Go versions in your project files.

Example usage:

	latestVersion, err := goautomate.GetLatestVersion()
	if err != nil {
	    log.Fatal(err)
	}

	err = goautomate.DownloadGo(latestVersion, "linux", "amd64")
	if err != nil {
	    log.Fatal(err)
	}

For more detailed information and advanced usage, refer to the README.md file.
*/
package goautomate
