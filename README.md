# GoAutomate🐿️

[![Release](https://github.com/Nicconike/goautomate/actions/workflows/release.yml/badge.svg)](https://github.com/Nicconike/goautomate/actions/workflows/release.yml)
[![Code Coverage](https://github.com/Nicconike/goautomate/actions/workflows/coverage.yml/badge.svg)](https://github.com/Nicconike/goautomate/actions/workflows/coverage.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nicconike/goautomate)](https://goreportcard.com/report/github.com/Nicconike/goautomate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nicconike/goautomate)
![GitHub Release](https://img.shields.io/github/v/release/nicconike/goautomate)
[![codecov](https://codecov.io/gh/Nicconike/goautomate/graph/badge.svg?token=MPIX1QLEYJ)](https://codecov.io/gh/Nicconike/goautomate)
[![Go Reference](https://pkg.go.dev/badge/github.com/Nicconike/goautomate.svg)](https://pkg.go.dev/github.com/Nicconike/goautomate)
![GitHub License](https://img.shields.io/github/license/nicconike/goautomate)
[![Visitor Badge](https://badges.pufler.dev/visits/nicconike/goautomate)](https://badges.pufler.dev)

goautomate is a Go tool that automates the process of checking and updating Go versions in your projects. It can detect the current Go version from various file types, compare it with the latest available version, and download the newest version if an update is available.

## Features

- Detect current Go version from various file types (Dockerfile, go.mod, JSON configs, etc.)
- Check for the latest available Go version
- Download the latest Go version if an update is available
- Checksum validation for downloaded Go versions to ensure integrity
- Support for different operating systems and architectures

## Installation

To install goautomate, use the following command:
```sh
go install github.com/Nicconike/goautomate/cmd/goautomate@latest
```

## Usage

### Basic Usage

```sh
goautomate -file <path-to-file> -os <target-os> -arch <target-arch>
```

This will check the specified file for the current Go version, compare it with the latest available version, and download the new version if an update is available.

> [!NOTE]
> If you don't specify the `os` and `arch` type, the tool will download the latest version for your current operating system and architecture.
>
> Minimum required Go version: 1.15

### Command-line Options

- `-file` or `-f`: Path to the file containing the current Go version
- `-version` or `-v`: Directly specify the current Go version
- `-os`: Target operating system (windows, linux, macOS[darwin])
- `-arch`: Target architecture (386[x86], amd64[x86-64], arm64, armv6l[armv6])

### Examples

1. Get version from a Dockerfile:
	```sh
	goautomate -f Dockerfile
	```
	![Dockerfile Example](https://github.com/Nicconike/goautomate/blob/master/assets/dockerfile_example.png)

2. Get version from go.mod:
	```sh
	goautomate -f go.mod
	```
	![Go Mod Example](https://github.com/Nicconike/goautomate/blob/master/assets/gomod_example.png)

3. Specify version directly:
	```sh
	goautomate -v 1.17
	```
	![Direct Example](https://github.com/Nicconike/goautomate/blob/master/assets/direct_example.png)

4. Download for a specific OS and architecture:
	```sh
	goautomate -f version.json -os linux -arch arm64
	```
	![JSON Example with OS](https://github.com/Nicconike/goautomate/blob/master/assets/json_example_os_arch.png)

## Supported File Types

`goautomate` can extract Go versions from various file types, including:

- Dockerfile
- go.mod
- JSON configuration files
- Plain text files with version information

The tool uses various regex patterns to detect Go versions, making it flexible for different project setups.

Missing any file types you expected to see? Let me know via [discussions](https://github.com/Nicconike/goautomate/discussions) or [discord server](https://discord.gg/UbetHfu).

## Contributing

Star⭐ and Fork🍴 the Repo to start with your feature request(or bug) and experiment with the project to implement whatever Idea you might have and sent the Pull Request through 🤙

Please refer [Contributing.md](https://github.com/Nicconike/Steam-Stats/blob/master/.github/CONTRIBUTING.md) to get to know how to contribute to this project.
And thank you for considering to contribute.

## License

[GPLv3 License](LICENSE)
