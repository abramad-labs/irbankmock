package sep

import (
	"regexp"
	"strings"
)

func SplitByDelimiters(s string) []string {
	re := regexp.MustCompile(`[|,;]`)
	parts := re.Split(s, -1)

	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func maskThirdQuarter(card string) string {
	if len(card) != 16 {
		return card
	}

	masked := card[:8] + strings.Repeat("*", 4) + card[12:]
	return masked
}
