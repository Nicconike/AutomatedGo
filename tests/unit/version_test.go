package tests

import (
	"testing"

	"github.com/Nicconike/goautomate/pkg"
)

func TestIsNewer(t *testing.T) {
	testCases := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		{"Newer major version", "go1.16.0", "go1.15.0", true},
		{"Newer minor version", "go1.15.2", "go1.15.1", true},
		{"Newer patch version", "go1.15.1", "go1.15.0", true},
		{"Same version", "go1.15.0", "go1.15.0", false},
		{"Older major version", "go1.14.0", "go1.15.0", false},
		{"Older minor version", "go1.15.1", "go1.15.2", false},
		{"Older patch version", "go1.15.0", "go1.15.1", false},
		{"Different length (newer)", "go1.15.1", "go1.15", true},
		{"Different length (older)", "go1.15", "go1.15.1", false},
		{"Without 'go' prefix", "1.16.0", "1.15.0", true},
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
