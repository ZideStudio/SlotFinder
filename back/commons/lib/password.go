package lib

import "regexp"

const PasswordMinLength = 8

var (
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasDigit   = regexp.MustCompile(`\d`)
	hasSpecial = regexp.MustCompile(`[^A-Za-z0-9]`)
)

// IsValidPassword validates that a password meets the following requirements:
// - At least 8 characters long
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one digit
// - Contains at least one special character (non-alphanumeric)
//
// This matches the frontend regex: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).+$/
func IsValidPassword(password string) bool {
	return len(password) >= PasswordMinLength && hasLower.MatchString(password) && hasUpper.MatchString(password) && hasDigit.MatchString(password) && hasSpecial.MatchString(password)
}
