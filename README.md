# AutomatedGo🐿️
[![Release](https://github.com/Nicconike/AutomatedGo/actions/workflows/release.yml/badge.svg)](https://github.com/Nicconike/AutomatedGo/actions/workflows/release.yml)
[![Publish Packages](https://github.com/Nicconike/AutomatedGo/actions/workflows/docker.yml/badge.svg)](https://github.com/Nicconike/AutomatedGo/actions/workflows/docker.yml)
[![CodeQL](https://github.com/Nicconike/AutomatedGo/actions/workflows/codeql.yml/badge.svg)](https://github.com/Nicconike/AutomatedGo/actions/workflows/codeql.yml)
[![Code Coverage](https://github.com/Nicconike/AutomatedGo/actions/workflows/coverage.yml/badge.svg)](https://github.com/Nicconike/AutomatedGo/actions/workflows/coverage.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nicconike/AutomatedGo)](https://goreportcard.com/report/github.com/Nicconike/AutomatedGo)
[![codecov](https://codecov.io/gh/Nicconike/AutomatedGo/graph/badge.svg?token=MPIX1QLEYJ)](https://codecov.io/gh/Nicconike/AutomatedGo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nicconike/AutomatedGo)
![GitHub Release](https://img.shields.io/github/v/release/nicconike/AutomatedGo)
![Docker Image Size](https://img.shields.io/docker/image-size/nicconike/automatedgo/master?sort=semver&logo=docker&label=Docker%20Image)
![Docker Pulls](https://img.shields.io/docker/pulls/nicconike/automatedgo?logo=docker&label=Docker%20Pulls)
[![Go Reference](https://pkg.go.dev/badge/github.com/Nicconike/AutomatedGo/v2.svg)](https://pkg.go.dev/github.com/Nicconike/AutomatedGo/v2)
![GitHub License](https://img.shields.io/github/license/nicconike/AutomatedGo)
[![wakatime](https://wakatime.com/badge/user/018e538b-3f55-4e8e-95fa-6c3225418eed/project/148b8322-28da-4cf4-85c2-bb20c2fe1295.svg)](https://wakatime.com/badge/user/018e538b-3f55-4e8e-95fa-6c3225418eed/project/148b8322-28da-4cf4-85c2-bb20c2fe1295)
[![Visitor Badge](https://badges.pufler.dev/visits/nicconike/AutomatedGo)](https://badges.pufler.dev)

**AutomatedGo** is a Go tool that automates the process of checking and updating Go versions in your projects. It can detect the current Go version from various file types, compare it with the latest available version, and download the newest version if an update is available.

## Features

- Detect current Go version from various file types (Dockerfile, go.mod, JSON configs, etc.)
- Check for the latest available Go version
- Download the latest Go version if an update is available
- Checksum validation for downloaded Go versions to ensure integrity
- Support for different operating systems and architectures

## Installation

To add **AutomatedGo** in your Go project, use the following command:
```sh
go get -u github.com/Nicconike/AutomatedGo/v2
```

To install **AutomatedGo** as a Go binary, use the following command:
```sh
go install github.com/Nicconike/AutomatedGo/v2/cmd/automatedgo@v2.1.0
```

## Usage

### Basic Usage

```sh
automatedgo -file <path-to-file> -os <target-os> -arch <target-arch>
```

This will check the specified file for the current Go version, compare it with the latest available version, and download the new version if an update is available.

> [!NOTE]
> If you don't specify the `os` and `arch` type, the tool will download the latest version for your current operating system and architecture.
>
> Minimum required Go version: 1.22

### Command-line Options

- `-file` or `-f`: Path to the file containing the current Go version
- `-version` or `-v`: Directly specify the current Go version
- `-os`: Target operating system (windows, linux, macOS[darwin])
- `-arch`: Target architecture (386[x86], amd64[x86-64], arm64, armv6l[armv6])

### Examples

1. Get version from a Dockerfile:
	```sh
	automatedgo -f Dockerfile
	```
	![Dockerfile Example](https://github.com/Nicconike/AutomatedGo/blob/master/assets/dockerfile_example.png)

2. Get version from go.mod:
	```sh
	automatedgo -f go.mod
	```
	![Go Mod Example](https://github.com/Nicconike/AutomatedGo/blob/master/assets/gomod_example.png)

3. Specify version directly:
	```sh
	automatedgo -v 1.18
	```
	![Direct Example](https://github.com/Nicconike/AutomatedGo/blob/master/assets/direct_example.png)

4. Download for a specific OS and architecture:
	```sh
	automatedgo -f version.json -os linux -arch arm64
	```
	![JSON Example with OS](https://github.com/Nicconike/AutomatedGo/blob/master/assets/json_example_os_arch.png)

> Also, checkout a real example in the [test-AutomatedGo](https://github.com/Nicconike/test-AutomatedGo) repository where this tool is used to check and update the Go version. And then upload the downloaded Go version to Github using Git LFS.

## Supported File Types

`AutomatedGo` can extract Go versions from various file types, including:

- Dockerfile
- go.mod
- JSON configuration files
- Plain text files with version information

The tool uses various regex patterns to detect Go versions, making it flexible for different project setups.

Missing any file types you expected to see? Let me know via [discussions](https://github.com/Nicconike/AutomatedGo/discussions) or [discord server](https://discord.gg/UbetHfu).

## Contributing

Star⭐ and Fork🍴 the Repo to start with your feature request(or bug) and experiment with the project to implement whatever Idea you might have and sent the Pull Request through 🤙

Please refer [Contributing.md](https://github.com/Nicconike/AutomatedGo/blob/master/.github/CONTRIBUTING.md) to get to know how to contribute to this project.
And thank you for considering to contribute.

## License

[GPLv3 License](LICENSE)
