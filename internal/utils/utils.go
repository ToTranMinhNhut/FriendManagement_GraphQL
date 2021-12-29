package utils

import "regexp"

const EmailRegex = `[_A-Za-z0-9-\+]+(\.[_A-Za-z0-9-]+)*@[A-Za-z0-9-]+(\.[A-Za-z0-9]+)*(\.[A-Za-z]{2,})`

func IsValidEmail(email string) (bool, error) {
	isValid, err := regexp.MatchString(EmailRegex, email)
	if err != nil || !isValid {
		return false, err
	}
	return true, nil
}

// Spit emails in a text
func GetMentionedEmailFromText(text string) []string {
	regex := regexp.MustCompile(EmailRegex)

	emailChain := regex.FindAllString(text, -1)
	email := make([]string, len(emailChain))
	for index, emailCharacter := range emailChain {
		email[index] = emailCharacter
	}
	return email
}
