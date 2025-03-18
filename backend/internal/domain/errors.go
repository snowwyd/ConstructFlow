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
	ErrRoleNotFound = errors.New("role not found")

	ErrRoleAlreadyExists = errors.New("role already exists")
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
)

var (
	ErrFileNotFound = errors.New("file not found")
)

var (
	ErrAccessDenied = errors.New("access denied")
)
