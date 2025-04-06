package domain

import "errors"

var (
	ErrInternal  = errors.New("internal error")
	ErrInvalidID = errors.New("invalid identifier")
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
	ErrFileNotFound                   = errors.New("file not found")
	ErrInvalidFileStatus              = errors.New("file is not in a draft state")
	ErrDirectoryContainsNonDraftFiles = errors.New("directory contains files with status other than 'draft'")
	ErrCannotDeleteNonDraftFile       = errors.New("cannot delete file with status other than 'draft'")

	ErrDirectoryNotFound = errors.New("directory not found")
)

var (
	ErrAccessDenied = errors.New("access denied")
)

var (
	ErrApprovalNotFound = errors.New("approval not found")
	ErrNoPermission     = errors.New("user has no permission to sign this approval")
)

var (
	ErrWorkflowNotFound = errors.New("workflow not found")
)
