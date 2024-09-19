package pkg

type VersionService struct {
	Downloader FileDownloader
	Remover    FileRemover
	Checksum   ChecksumCalculator
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

func (v *VersionService) DownloadGo(version, targetOS, arch, path string) error {
	return DownloadGo(version, targetOS, arch, path, v.Downloader, v.Remover, v.Checksum)
}
