package lib

import "regexp"

const (
	PasswordMinLength = 8
)

var (
	// Regular expressions to validate password requirements
	// Since Go's regexp doesn't support lookahead assertions, we check each requirement separately
	hasLowerRegex   = regexp.MustCompile(`[a-z]`)
	hasUpperRegex   = regexp.MustCompile(`[A-Z]`)
	hasDigitRegex   = regexp.MustCompile(`\d`)
	hasSpecialRegex = regexp.MustCompile(`[^A-Za-z0-9]`)
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
	if len(password) < PasswordMinLength {
		return false
	}

	return hasLowerRegex.MatchString(password) &&
		hasUpperRegex.MatchString(password) &&
		hasDigitRegex.MatchString(password) &&
		hasSpecialRegex.MatchString(password)
}
