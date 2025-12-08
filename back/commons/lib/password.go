package lib

import (
	"unicode"
)

const (
	PasswordMinLength = 8
)

// IsValidPassword validates that a password meets the following requirements:
// - At least 8 characters long
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one digit
// - Contains at least one special character (non-alphanumeric)
func IsValidPassword(password string) bool {
	if len(password) < PasswordMinLength {
		return false
	}

	var (
		hasLower   = false
		hasUpper   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

