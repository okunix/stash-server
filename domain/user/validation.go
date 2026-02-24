package user

import (
	"errors"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{2,20}$`)
	passwordRegex = regexp.MustCompile(`^.{8,}$`)
)

var (
	ErrInvalidUsername = errors.New("Invalid username provided")
	ErrInvalidPassword = errors.New("Invalid password provided")
	ErrNoPasswordHash  = errors.New("password hash must be provided")
)

func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func ValidatePassword(password string) error {
	if !passwordRegex.MatchString(password) {
		return ErrInvalidPassword
	}
	return nil
}

func ValidatePasswordHash(passwordHash string) error {
	if len(passwordHash) < 1 {
		return ErrNoPasswordHash
	}
	return nil
}
