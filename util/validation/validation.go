package validation

import (
	"regexp"

	"../log"
)

// EmailRegex ...
const EmailRegex = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`

// EmailValidation returns valid of email
func EmailValidation(email string) bool {
	exp, err := regexp.Compile(EmailRegex)
	if regexpCompiled := log.CheckError(err); regexpCompiled {
		if exp.MatchString(email) {
			return true
		}
	}
	return false
}
