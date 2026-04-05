package lib

import (
	"regexp"
	"strings"
)

// IsHexa checks if the string is in hexadecimal color format
func IsHexa(s string) bool {
	if len(s) != 7 {
		return false
	}

	matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, s)
	return matched
}

// BoolToString converts a boolean value to its string representation
func BoolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// Capitalize returns the string with its first rune uppercased
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = []rune(strings.ToUpper(string(r[0])))[0]
	return string(r)
}
