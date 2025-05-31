package ton

import (
	"fmt"
	"strconv"
	"strings"
)

func averagePercent(percentStrings []string) (float64, error) {
	var sum float64
	for _, p := range percentStrings {
		// Remove trailing '%' sign
		trimmed := strings.TrimSuffix(p, "%")
		// Parse to float64
		val, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse %q: %w", p, err)
		}
		sum += val
	}
	if len(percentStrings) == 0 {
		return 0, nil // avoid division by zero
	}
	return sum / float64(len(percentStrings)), nil
}

func compareStrings(s1, s2 string) string {
	// Determine the minimum length to avoid index out of range
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	// Count matching characters at the same positions
	matchCount := 0
	for i := 0; i < minLen; i++ {
		if s1[i] == s2[i] {
			matchCount++
		}
	}

	// Use the length of the longer string as the denominator
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	// Calculate the percentage
	percent := (float64(matchCount) / float64(maxLen)) * 100

	// Format with two decimal places
	return fmt.Sprintf("%.2f%%", percent)
}
