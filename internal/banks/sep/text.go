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
