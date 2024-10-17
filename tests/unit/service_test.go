package tests

import (
	"bytes"
	"testing"

	"github.com/Nicconike/AutomatedGo/v2/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockFileDownloader struct {
	mock.Mock
}

func (m *MockFileDownloader) Download(url, filename string) error {
	args := m.Called(url, filename)
	return args.Error(0)
}

type MockFileRemover struct {
	mock.Mock
}

func (m *MockFileRemover) Remove(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

// func TestVersionServiceGetCurrentVersion(t *testing.T) {
// 	vs := &pkg.VersionService{}
// 	_, err := vs.GetCurrentVersion("testfile", "1.0.0")
// 	assert.NoError(t, err)
// }

func TestVersionServiceGetLatestVersion(t *testing.T) {
	vs := &pkg.VersionService{}
	_, err := vs.GetLatestVersion()
	assert.NoError(t, err)
}

func TestVersionServiceIsNewer(t *testing.T) {
	vs := &pkg.VersionService{}
	result := vs.IsNewer("1.2.0", "1.1.0")
	assert.True(t, result)
}

func TestVersionServiceDownloadGo(t *testing.T) {
	mockDownloader := new(MockFileDownloader)
	mockRemover := new(MockFileRemover)
	mockChecksum := new(MockChecksumCalculator)

	vs := &pkg.VersionService{
		Downloader: mockDownloader,
		Remover:    mockRemover,
		Checksum:   mockChecksum,
	}

	mockChecksum.On("GetOfficialChecksum", mock.Anything).Return("checksum", nil)
	mockDownloader.On("Download", mock.Anything, mock.Anything).Return(nil)
	mockChecksum.On("Calculate", mock.Anything).Return("checksum", nil)

	input := bytes.NewBufferString("")
	output := &bytes.Buffer{}

	err := vs.DownloadGo("1.16.5", "linux", "amd64", "/tmp", input, output)
	assert.NoError(t, err)

	mockDownloader.AssertExpectations(t)
	mockRemover.AssertExpectations(t)
	mockChecksum.AssertExpectations(t)
}
