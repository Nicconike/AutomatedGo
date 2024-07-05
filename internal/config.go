package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
)

func GetCurrentVersion(filePath, directVersion string) (string, error) {
	if filePath != "" {
		return readVersionFromFile(filePath)
	}
	if directVersion != "" {
		return directVersion, nil
	}
	return "", errors.New("no version input provided")
}

func readVersionFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	version := extractGoVersion(string(content))
	if version != "" {
		return version, nil
	}

	return "", errors.New("unable to extract Go version from file")
}

func extractGoVersion(content string) string {
	// Common patterns for Go version
	patterns := []string{
		`(?i)go\s*version\s*[:=]?\s*["']?(\d+\.\d+(\.\d+)?)["']?`,
		`(?i)go_version\s*[:=]?\s*["']?(\d+\.\d+(\.\d+)?)["']?`,
		`(?i)golang_version\s*[:=]?\s*["']?(\d+\.\d+(\.\d+)?)["']?`,
		`(?i)go(\d+\.\d+(\.\d+)?)`,
		`(?i)FROM\s+golang:(\d+\.\d+(\.\d+)?)`,
		`(?i)ARG\s+GO_VERSION=(\d+\.\d+(\.\d+)?)`,
		`(?i)ENV\s+GO_VERSION=(\d+\.\d+(\.\d+)?)`,
	}

	// Check for JSON format
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err == nil {
		for _, key := range []string{"go_version", "goVersion", "golang_version", "golangVersion", "GO_VERSION"} {
			if version, ok := jsonData[key].(string); ok {
				return version
			}
		}
	}

	// Check for go.mod file
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "go ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}

	// Check for other formats using regex
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(content, -1)
		if len(matches) > 0 {
			return matches[len(matches)-1][1]
		}
	}

	return ""
}
