package pkg

import (
	"strconv"
	"strings"
)

func IsNewer(latest, current string) bool {
	latestParts := strings.Split(strings.TrimPrefix(latest, "go"), ".")
	currentParts := strings.Split(strings.TrimPrefix(current, "go"), ".")

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		latestNum, _ := strconv.Atoi(latestParts[i])
		currentNum, _ := strconv.Atoi(currentParts[i])
		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	return len(latestParts) > len(currentParts)
}
