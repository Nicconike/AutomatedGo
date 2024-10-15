package pkg

import (
	"io"
)

type VersionService struct {
	Downloader FileDownloader
	Remover    FileRemover
	Checksum   ChecksumCalculator
	Input      io.Reader
	Output     io.Writer
}

func (v *VersionService) GetCurrentVersion(versionFile, currentVersion string) (string, error) {
	return GetCurrentVersion(versionFile, currentVersion)
}

func (v *VersionService) GetLatestVersion() (string, error) {
	return GetLatestVersion()
}

func (v *VersionService) IsNewer(latestVersion, currentVersion string) bool {
	return IsNewer(latestVersion, currentVersion)
}

func (v *VersionService) DownloadGo(version, targetOS, arch, path string, input io.Reader, output io.Writer) error {
	config := DownloadConfig{
		Version:    version,
		TargetOS:   targetOS,
		Arch:       arch,
		Path:       path,
		Downloader: v.Downloader,
		Remover:    v.Remover,
		Checksum:   v.Checksum,
		Input:      input,
		Output:     output,
	}
	return DownloadGo(config)
}
