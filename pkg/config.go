package pkg

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
		return ReadVersionFromFile(filePath)
	}
	if directVersion != "" {
		return directVersion, nil
	}
	return "", errors.New("no version input provided")
}

func ReadVersionFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	version := ExtractGoVersion(string(content))
	if version != "" {
		return version, nil
	}

	return "", errors.New("unable to extract Go version from file")
}

func ExtractGoVersion(content string) string {
	// Common patterns for Go version
	patterns := []string{
		`(?i)(?:go|golang|go_version|golang_version)(?:\s*version)?[:=]?\s*v?(\d+\.\d+(?:\.\d+)?)`,
		`(?i)FROM\s+golang:(\d+\.\d+(?:\.\d+)?)`,
		`(?i)ARG\s+GO_VERSION=(\d+\.\d+(?:\.\d+)?)`,
		`(?i)ENV\s+GO_VERSION=(\d+\.\d+(?:\.\d+)?)`,
		`(\d+\.\d+(?:\.\d+)?)`,
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
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}
