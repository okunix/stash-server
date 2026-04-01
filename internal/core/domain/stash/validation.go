package stash

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_.\/]{1,255}$`)
	descRegex = regexp.MustCompile(`^.{0,1000}$`)
)

var (
	ErrInvalidName     = errors.New("Invalid name provided")
	ErrInvalidDesc     = errors.New("Invalid description provided")
	ErrNoMasterKeyHash = errors.New("Master key hash must be provided")
	ErrNoMasterKeySalt = errors.New("Master key salt must be provided")
)

func ValidateName(name string) error {
	if len(name) < 1 {
		return errors.New("Name is too short")
	}
	if len(name) > 255 {
		return errors.New("Name is too long")
	}
	if _, err := uuid.Parse(name); err == nil {
		return errors.New("Name must not be a valid UUID")
	}
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}

func ValidateDescription(description *string) error {
	if description == nil {
		return nil
	}
	if len(*description) > 1000 {
		return errors.New("Description is too long")
	}
	if !descRegex.MatchString(*description) {
		return ErrInvalidDesc
	}
	return nil
}

func ValidateMasterKeyHash(hash string) error {
	if len(hash) < 1 {
		return ErrNoMasterKeyHash
	}
	return nil
}

func ValidateMasterKeySalt(salt string) error {
	if len(salt) < 1 {
		return ErrNoMasterKeySalt
	}
	return nil
}

func ValidatePassword(key string) error {
	if len(key) < 10 {
		return errors.New("Password is too short")
	}
	return nil
}
