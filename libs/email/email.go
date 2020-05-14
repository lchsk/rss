package email

import (
	"regexp"
)

const strEmailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var regexEmail = regexp.MustCompile(strEmailRegex)

func IsEmailValid(email string) bool {
	if len(email) > 254 || !regexEmail.MatchString(email) {
		return false
	}

	return true
}
