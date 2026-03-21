package secret

import (
	"errors"
	"regexp"
)

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_+. ]+$`)
)

func ValidateEntryName(name string) error {
	if !nameRegex.MatchString(name) {
		return errors.New("entry name is not valid")
	}
	return nil
}
