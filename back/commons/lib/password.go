package lib

import "regexp"

const PasswordMinLength = 8

// IsValidPassword validates password security requirements
func IsValidPassword(password string) bool {
	return len(password) >= PasswordMinLength && 
		regexp.MustCompile(`[a-z]`).MatchString(password) && 
		regexp.MustCompile(`[A-Z]`).MatchString(password) && 
		regexp.MustCompile(`\d`).MatchString(password) && 
		regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
}
