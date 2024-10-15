package tests

import (
	"testing"

	"github.com/Nicconike/AutomatedGo/v2/pkg"
)

func TestIsNewer(t *testing.T) {
	testCases := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		{"Newer major version", "go2.0.0", "go1.16.0", true},
		{"Newer minor version", "go1.16.5", "go1.15.10", true},
		{"Newer patch version", "go1.15.10", "go1.15.9", true},
		{"Same version with different prefix", "go1.15.0", "1.15.0", false},
		{"Older major version", "go1.14.0", "go2.0.0", false},
		{"Older minor version with same major", "go1.15.9", "go1.16.1", false},
		{"Older patch version with same major and minor", "go1.17.0", "go1.17.1", false},
		{"Different length (newer)", "go1.17.10", "go1.17", true},
		{"Different length (older)", "go1.18", "go1.18.10", false},
		{"Without 'go' prefix newer major", "2.0.0", "1.16.0", true},
		{"Without 'go' prefix older minor", "1.14.5", "1.14.6", false},
		{"Without 'go' prefix same version different format", "1.16.0", "1.16", true},
		{"Empty latest version string", "", "go1.20.0", false},
		{"Empty current version string", "go1.16.0", "", true},
		{"Both versions empty strings", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pkg.IsNewer(tc.latest, tc.current)
			if result != tc.expected {
				t.Errorf("IsNewer(%q, %q) = %v; want %v", tc.latest, tc.current, result, tc.expected)
			}
		})
	}
}
