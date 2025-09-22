package lib

import "regexp"

const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

func IsValidEmail(email string) bool {
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}
