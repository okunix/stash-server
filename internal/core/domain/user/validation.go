package user

import (
	"errors"
	"fmt"
	"regexp"
	"slices"

	"github.com/google/uuid"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{2,20}$`)
	passwordRegex = regexp.MustCompile(`^.{10,}$`)
)

const (
	passwordMinLength = 10
)

var (
	ErrInvalidUsername = errors.New("Invalid username provided")
	ErrInvalidPassword = errors.New("Invalid password provided")
	ErrNoPasswordHash  = errors.New("Password hash must be provided")
)

var (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func ValidateUsername(username string) error {
	if len(username) < 2 {
		return errors.New("Username is too short")
	}
	if len(username) > 40 {
		return errors.New("Username is too long")
	}
	if _, err := uuid.Parse(username); err == nil {
		return errors.New("Username must not be valid uuid")
	}
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < passwordMinLength {
		return errors.New("Password is too short")
	}
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

func ValidateUserRole(role string) error {
	roles := []string{RoleAdmin, RoleUser}
	if !slices.Contains(roles, role) {
		return fmt.Errorf("Unknown role specified: %s", role)
	}
	return nil
}
