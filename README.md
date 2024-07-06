# GoAutomate

[![Release](https://github.com/Nicconike/goautomate/actions/workflows/release.yml/badge.svg)](https://github.com/Nicconike/goautomate/actions/workflows/release.yml)
[![Code Coverage](https://github.com/Nicconike/goautomate/actions/workflows/coverage.yml/badge.svg)](https://github.com/Nicconike/goautomate/actions/workflows/coverage.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nicconike/goautomate)
![GitHub Release](https://img.shields.io/github/v/release/nicconike/goautomate)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nicconike/goautomate)](https://goreportcard.com/report/github.com/Nicconike/goautomate)
[![codecov](https://codecov.io/gh/Nicconike/goautomate/graph/badge.svg?token=MPIX1QLEYJ)](https://codecov.io/gh/Nicconike/goautomate)
[![Go Reference](https://pkg.go.dev/badge/github.com/Nicconike/goautomate.svg)](https://pkg.go.dev/github.com/Nicconike/goautomate)
![GitHub License](https://img.shields.io/github/license/nicconike/goautomate)

goautomate is a Go tool that automates the process of checking and updating Go versions in your projects. It can detect the current Go version from various file types, compare it with the latest available version, and download the newest version if an update is available.

## Features

- Detect current Go version from various file types (Dockerfile, go.mod, JSON configs, etc.)
- Check for the latest available Go version
- Download the latest Go version if an update is available
- Support for different operating systems and architectures
- Progress bar display during download

## Installation

To install goautomate, use the following command:
```sh
go get -u github.com/Nicconike/goautomate
```

## Usage

### Basic Usage

```sh
goautomate -file <path-to-file>
```

This will check the specified file for the current Go version, compare it with the latest available version, and download the new version if an update is available.

### Command-line Options

- `-file`: Path to the file containing the current Go version
- `-version` or `-v`: Directly specify the current Go version
- `-os`: Target operating system (windows, linux, darwin)
- `-arch`: Target architecture (386, amd64, arm64, armv6l)

### Examples

1. Get version from a Dockerfile:
	```sh
	goautomate -file ./Dockerfile
	```

2. Get version from go.mod:
	```sh
	goautomate -file ./go.mod
	```

3. Specify version directly:
	```sh
	goautomate -version 1.16.5
	```

4. Download for a specific OS and architecture:
	```sh
	goautomate -file ./go.mod -os linux -arch arm64
	```

## Supported File Types

goautomate can extract Go versions from various file types, including:

- Dockerfile
- go.mod
- JSON configuration files
- Plain text files with version information

The tool uses various patterns to detect Go versions, making it flexible for different project setups.

## Contributing

Contributions to goautomate are welcome! Please feel free to submit a Pull Request.

## License

[GPLv3 License](LICENSE)
