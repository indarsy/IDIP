package utils

import "errors"

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrDuplicateRecord  = errors.New("duplicate record")
	ErrStorageNotFound  = errors.New("storage not found")
	ErrInvalidStatus    = errors.New("invalid status")
	ErrOperationFailed  = errors.New("operation failed")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrFileNotFound     = errors.New("file not found")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrFileTooLarge     = errors.New("file too large")
)
