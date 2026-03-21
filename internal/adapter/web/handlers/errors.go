package handlers

import "errors"

var (
	errUserNotFound  = errors.New("user not found")
	errStashNotFound = errors.New("stash not found")
)
