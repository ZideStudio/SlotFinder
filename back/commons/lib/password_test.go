package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		// Valid passwords
		{
			name:     "Valid password with all requirements",
			password: "Test123!",
			expected: true,
		},
		{
			name:     "Valid password with special character @",
			password: "Password1@",
			expected: true,
		},
		{
			name:     "Valid password with special character #",
			password: "MyPass123#",
			expected: true,
		},
		{
			name:     "Valid password with special character $",
			password: "Secure$123",
			expected: true,
		},
		{
			name:     "Valid password with multiple special characters",
			password: "P@ssw0rd!",
			expected: true,
		},
		{
			name:     "Valid long password",
			password: "VerySecure123!@#",
			expected: true,
		},

		// Invalid passwords - too short
		{
			name:     "Password too short - 7 characters",
			password: "Test12!",
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			expected: false,
		},
		{
			name:     "Password with only 1 character",
			password: "A",
			expected: false,
		},

		// Invalid passwords - missing required character types
		{
			name:     "Missing uppercase letter",
			password: "password123!",
			expected: false,
		},
		{
			name:     "Missing lowercase letter",
			password: "PASSWORD123!",
			expected: false,
		},
		{
			name:     "Missing digit",
			password: "Password!",
			expected: false,
		},
		{
			name:     "Missing special character",
			password: "Password123",
			expected: false,
		},
		{
			name:     "Only lowercase letters",
			password: "abcdefgh",
			expected: false,
		},
		{
			name:     "Only uppercase letters",
			password: "ABCDEFGH",
			expected: false,
		},
		{
			name:     "Only digits",
			password: "12345678",
			expected: false,
		},
		{
			name:     "Only special characters",
			password: "!@#$%^&*",
			expected: false,
		},
		{
			name:     "Lowercase and uppercase only",
			password: "Abcdefgh",
			expected: false,
		},
		{
			name:     "Lowercase and digits only",
			password: "abcd1234",
			expected: false,
		},
		{
			name:     "Lowercase and special chars only",
			password: "abcd!@#$",
			expected: false,
		},
		{
			name:     "Uppercase and digits only",
			password: "ABCD1234",
			expected: false,
		},
		{
			name:     "Uppercase and special chars only",
			password: "ABCD!@#$",
			expected: false,
		},
		{
			name:     "Digits and special chars only",
			password: "1234!@#$",
			expected: false,
		},

		// Edge cases
		{
			name:     "Exactly 8 characters with all requirements",
			password: "Aa1!bcde",
			expected: true,
		},
		{
			name:     "Password with spaces (valid special char)",
			password: "Pass 123!",
			expected: true,
		},
		{
			name:     "Password with unicode characters",
			password: "Pässw0rd!",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidPassword(tc.password)
			assert.Equal(t, tc.expected, result, "Password validation failed for: %s", tc.password)
		})
	}
}

func TestPasswordMinLength(t *testing.T) {
	assert.Equal(t, 8, PasswordMinLength, "PasswordMinLength should be 8")
}
