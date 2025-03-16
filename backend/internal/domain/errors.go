package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
)

var (
	ErrUserNotFound = errors.New("user not found")

	ErrUserAlreadyExists = errors.New("user already exists")
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
)
