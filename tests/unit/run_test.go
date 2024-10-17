package tests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Nicconike/AutomatedGo/v2/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVersionChecker struct {
	mock.Mock
}

func (m *MockVersionChecker) GetLatestVersion() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockVersionChecker) GetCurrentVersion(versionFile, currentVersion string) (string, error) {
	args := m.Called(versionFile, currentVersion)
	return args.String(0), args.Error(1)
}

func (m *MockVersionChecker) IsNewer(latestVersion, currentVersion string) bool {
	args := m.Called(latestVersion, currentVersion)
	return args.Bool(0)
}

func (m *MockVersionChecker) DownloadGo(version, targetOS, arch, path string, input io.Reader, output io.Writer) error {
	args := m.Called(version, targetOS, arch, path, input, output)
	return args.Error(0)
}

// Helper function to simulate user input and capture output
func simulateGetDownloadPath(input string) (string, string) {
	in := strings.NewReader(input)
	out := new(bytes.Buffer)
	path := pkg.GetDownloadPath(in, out)
	return path, out.String()
}

func TestGetDownloadPath(t *testing.T) {
	t.Run("Cancel input", func(t *testing.T) {
		path, output := simulateGetDownloadPath("cancel\n")
		assert.Equal(t, "", path)
		assert.Contains(t, output, "Enter the path where you want to download the file")
	})

	t.Run("Empty input uses current directory", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		path, output := simulateGetDownloadPath("\n")
		assert.Equal(t, currentDir, path)
		assert.Contains(t, output, fmt.Sprintf("Using current directory: %s", currentDir))
	})

	// t.Run("Non-existent path", func(t *testing.T) {
	// 	nonExistentPath := "/this/path/does/not/exist"
	// 	path, output := simulateGetDownloadPath(nonExistentPath + "\n")
	// 	assert.Equal(t, "", path)
	// 	assert.Contains(t, output, "Specified path does not exist. Please try again.")
	// })

	t.Run("Valid path", func(t *testing.T) {
		tempDir := os.TempDir()
		path, output := simulateGetDownloadPath(tempDir + "\n")
		assert.Equal(t, tempDir, path)
		assert.NotContains(t, output, "Specified path does not exist. Please try again")
	})
}

func TestRun(t *testing.T) {
	const version = "version.txt"
	tests := []struct {
		name           string
		versionFile    string
		currentVersion string
		targetOS       string
		targetArch     string
		input          string
		expectedOutput string
		expectedError  error
		mockSetup      func(m *MockVersionChecker)
	}{
		{
			name:          "No version specified",
			expectedError: errors.New("error: Either -file (-f) or -version (-v) must be specified"),
			mockSetup: func(m *MockVersionChecker) {
				// No specific setup required for this test case
			},
		},
		{
			name:           "Latest version available, user agrees to download",
			versionFile:    version,
			currentVersion: "",
			targetOS:       "linux",
			targetArch:     "amd64",
			input:          "yes\n\n",
			expectedOutput: "Current version: 1.0.0\n" +
				"Latest version: 1.1.0\n" +
				"A newer version is available\n" +
				"Do you want to download the latest version? (yes/no): " +
				"Enter the path where you want to download the file (press Enter for current directory, or 'cancel' to abort): " +
				"Using current directory: {{.CurrentDir}}\n" +
				"1.1.0 has been downloaded to {{.CurrentDir}}\n",
			expectedError: nil,
			mockSetup: func(m *MockVersionChecker) {
				m.On("GetCurrentVersion", version, "").Return("1.0.0", nil)
				m.On("GetLatestVersion").Return("1.1.0", nil)
				m.On("IsNewer", "1.1.0", "1.0.0").Return(true)
				m.On("DownloadGo", "1.1.0", "linux", "amd64", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:           "Latest version available, user declines to download",
			versionFile:    version,
			currentVersion: "",
			targetOS:       "linux",
			targetArch:     "amd64",
			input:          "no\n",
			expectedOutput: "Current version: 1.0.0\nLatest version: 1.1.0\nA newer version is available\nDo you want to download the latest version? (yes/no): Download aborted by user\n",
			expectedError:  nil,
			mockSetup: func(m *MockVersionChecker) {
				m.On("GetCurrentVersion", version, "").Return("1.0.0", nil)
				m.On("GetLatestVersion").Return("1.1.0", nil)
				m.On("IsNewer", "1.1.0", "1.0.0").Return(true)
			},
		},
		{
			name:           "No new version available",
			versionFile:    version,
			currentVersion: "",
			targetOS:       "linux",
			targetArch:     "amd64",
			input:          "",
			expectedOutput: "Current version: 1.1.0\nLatest version: 1.1.0\nYou have the latest version\n",
			expectedError:  nil,
			mockSetup: func(m *MockVersionChecker) {
				m.On("GetCurrentVersion", version, "").Return("1.1.0", nil)
				m.On("GetLatestVersion").Return("1.1.0", nil)
				m.On("IsNewer", "1.1.0", "1.1.0").Return(false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVersionChecker)
			tt.mockSetup(mockService)

			input := bytes.NewBufferString(tt.input)
			output := new(bytes.Buffer)

			err := pkg.Run(mockService, tt.versionFile, tt.currentVersion, tt.targetOS, tt.targetArch, input, output)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				currentDir, _ := os.Getwd()
				expectedOutput := strings.ReplaceAll(tt.expectedOutput, "{{.CurrentDir}}", currentDir)
				assert.Equal(t, expectedOutput, output.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}
