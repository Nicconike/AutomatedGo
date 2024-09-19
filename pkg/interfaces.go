package pkg

type VersionChecker interface {
	GetCurrentVersion(versionFile, currentVersion string) (string, error)
	GetLatestVersion() (string, error)
	IsNewer(latestVersion, currentVersion string) bool
	DownloadGo(version, os, arch, path string) error
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
