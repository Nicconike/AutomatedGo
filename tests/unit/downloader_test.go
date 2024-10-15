package tests

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Nicconike/AutomatedGo/v2/pkg"
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
	const checksum = "checksum-value"
	tests := []struct {
		name          string
		config        pkg.DownloadConfig
		setupMocks    func(d *MockDownloader, r *MockRemover, c *MockChecksumCalculator)
		expectedError error
	}{
		{
			name: "Successful download and checksum verification",
			config: pkg.DownloadConfig{
				Version:  "1.16.5",
				TargetOS: "linux",
				Arch:     "amd64",
				Path:     "",
				Input:    strings.NewReader(""),
				Output:   &bytes.Buffer{},
			},
			setupMocks: func(d *MockDownloader, r *MockRemover, c *MockChecksumCalculator) {
				c.On("GetOfficialChecksum", mock.Anything).Return(checksum, nil)
				d.On("Download", mock.Anything, mock.Anything).Return(nil)
				c.On("Calculate", mock.Anything).Return(checksum, nil)
			},
			expectedError: nil,
		},
		{
			name: "Download with user input for OS and Arch",
			config: pkg.DownloadConfig{
				Version: "1.16.5",
				Path:    "",
				Input:   strings.NewReader("linux\namd64\n"),
				Output:  &bytes.Buffer{},
			},
			setupMocks: func(d *MockDownloader, r *MockRemover, c *MockChecksumCalculator) {
				c.On("GetOfficialChecksum", mock.Anything).Return(checksum, nil)
				d.On("Download", mock.Anything, mock.Anything).Return(nil)
				c.On("Calculate", mock.Anything).Return(checksum, nil)
			},
			expectedError: nil,
		},
		{
			name: "Checksum mismatch",
			config: pkg.DownloadConfig{
				Version:  "1.16.5",
				TargetOS: "linux",
				Arch:     "amd64",
				Path:     "",
				Input:    strings.NewReader(""),
				Output:   &bytes.Buffer{},
			},
			setupMocks: func(d *MockDownloader, r *MockRemover, c *MockChecksumCalculator) {
				c.On("GetOfficialChecksum", mock.Anything).Return("correct-checksum", nil)
				d.On("Download", mock.Anything, mock.Anything).Return(nil)
				c.On("Calculate", mock.Anything).Return("incorrect-checksum", nil)
				r.On("Remove", mock.Anything).Return(nil)
			},
			expectedError: errors.New("Checksum mismatch: expected correct-checksum, got incorrect-checksum for file go1.16.5.linux-amd64.tar.gz"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDownloader := new(MockDownloader)
			mockRemover := new(MockRemover)
			mockChecksumCalculator := new(MockChecksumCalculator)

			tt.setupMocks(mockDownloader, mockRemover, mockChecksumCalculator)

			tt.config.Downloader = mockDownloader
			tt.config.Remover = mockRemover
			tt.config.Checksum = mockChecksumCalculator

			err := pkg.DownloadGo(tt.config)

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
