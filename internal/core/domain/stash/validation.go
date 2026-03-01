package stash

import (
	"errors"
	"regexp"
)

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_.\/]{1,255}$`)
	descRegex = regexp.MustCompile(`^.{0,1000}$`)
)

var (
	ErrInvalidName     = errors.New("invalid name provided")
	ErrInvalidDesc     = errors.New("invalid description provided")
	ErrNoMasterKeyHash = errors.New("master key hash must be provided")
	ErrNoMasterKeySalt = errors.New("master key salt must be provided")
)

func ValidateName(name string) error {
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}

func ValidateDescription(description string) error {
	if !descRegex.MatchString(description) {
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
