package validator

import (
	"regexp"
	"unicode/utf8"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	return utf8.RuneCountInString(password) >= 8
}

func ValidateNickname(nickname string) bool {
	length := utf8.RuneCountInString(nickname)
	return length >= 2 && length <= 20
}
