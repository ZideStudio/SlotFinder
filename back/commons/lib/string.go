package lib

import (
	"regexp"
)

// isHexa checks if the string is in hexadecimal color format
func IsHexa(s string) bool {
	if len(s) != 7 {
		return false
	}

	matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, s)
	return matched
}
