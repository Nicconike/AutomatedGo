package tests

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/Nicconike/AutomatedGo/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDownloader is a mock implementation of the FileDownloader interface
type MockDownloader struct {
	mock.Mock
}

func (m *MockDownloader) Download(url, filename string) error {
	args := m.Called(url, filename)
	return args.Error(0)
}

// MockRemover is a mock implementation of the FileRemover interface
type MockRemover struct {
	mock.Mock
}

func (m *MockRemover) Remove(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

// MockChecksumCalculator is a mock implementation of the ChecksumCalculator interface
type MockChecksumCalculator struct {
	mock.Mock
}

func (m *MockChecksumCalculator) GetOfficialChecksum(filename string) (string, error) {
	args := m.Called(filename)
	return args.String(0), args.Error(1)
}

func (m *MockChecksumCalculator) Calculate(filename string) (string, error) {
	args := m.Called(filename)
	return args.String(0), args.Error(1)
}

func TestDownloadGo(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		targetOS      string
		arch          string
		path          string
		setupMocks    func(d *MockDownloader, r *MockRemover, c *MockChecksumCalculator)
		expectedError error
	}{
		{
			name:     "Successful download and checksum verification",
			version:  "go1.16.5",
			targetOS: runtime.GOOS,
			arch:     runtime.GOARCH,
			path:     "/tmp",
			setupMocks: func(d *MockDownloader, r *MockRemover, c *MockChecksumCalculator) {
				var extension string
				if runtime.GOOS == "windows" {
					extension = "zip"
				} else {
					extension = "tar.gz"
				}
				filename := fmt.Sprintf("/tmp/go1.16.5.%s-%s.%s", runtime.GOOS, runtime.GOARCH, extension)

				c.On("GetOfficialChecksum", filename).Return("checksum-value", nil)
				d.On("Download", fmt.Sprintf(pkg.DownloadURLFormat, "1.16.5", runtime.GOOS, runtime.GOARCH, extension), filename).Return(nil)
				c.On("Calculate", filename).Return("checksum-value", nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDownloader := new(MockDownloader)
			mockRemover := new(MockRemover)
			mockChecksumCalculator := new(MockChecksumCalculator)

			tt.setupMocks(mockDownloader, mockRemover, mockChecksumCalculator)

			err := pkg.DownloadGo(tt.version, tt.targetOS, tt.arch, tt.path, mockDownloader, mockRemover, mockChecksumCalculator)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDownloader.AssertExpectations(t)
			mockRemover.AssertExpectations(t)
			mockChecksumCalculator.AssertExpectations(t)
		})
	}
}
