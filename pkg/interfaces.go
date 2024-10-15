package pkg

import "io"

type VersionChecker interface {
	GetLatestVersion() (string, error)
	GetCurrentVersion(versionFile, currentVersion string) (string, error)
	IsNewer(latestVersion, currentVersion string) bool
	DownloadGo(version, targetOS, arch, path string, input io.Reader, output io.Writer) error
}

type FileDownloader interface {
	Download(url, filename string) error
}

type FileRemover interface {
	Remove(filename string) error
}

type ChecksumCalculator interface {
	Calculate(filename string) (string, error)
	GetOfficialChecksum(filename string) (string, error)
}
