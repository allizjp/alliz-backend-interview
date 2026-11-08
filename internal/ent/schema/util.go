package schema

import (
	"net/mail"
	"strings"
)

func notBlank(errBlank error) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errBlank
		}
		return nil
	}
}

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return ErrEventRequestEmailInvalid
	}
	return nil
}
