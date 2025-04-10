package domain

import "errors"

var (
	ErrInternal  = errors.New("internal error")
	ErrInvalidID = errors.New("invalid identifier")
)

var (
	ErrFileNotFound                   = errors.New("file not found")
	ErrEmptyFile                      = errors.New("empty file in MinIO")
	ErrInvalidFileStatus              = errors.New("file is not in a draft state")
	ErrDirectoryContainsNonDraftFiles = errors.New("directory contains files with status other than 'draft'")
	ErrCannotDeleteNonDraftFile       = errors.New("cannot delete file with status other than 'draft'")

	ErrDirectoryNotFound      = errors.New("directory not found")
	ErrDirectoryAlreadyExists = errors.New("directory already exists")
	ErrFileAlreadyExists      = errors.New("file already exists")
)

var (
	ErrAccessDenied = errors.New("access denied")
)

var (
	ErrApprovalNotFound = errors.New("approval not found")
	ErrNoPermission     = errors.New("user has no permission to sign this approval")
)

var (
	ErrUnsupportedFormat = errors.New("unsupported file format")
	ErrConversionFailed  = errors.New("conversion failed")
)

var (
	ErrNoRelationsFound = errors.New("no user relations found")
)
