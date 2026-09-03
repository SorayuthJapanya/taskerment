package domain

import "errors"

var (
	ErrNotFound          = errors.New("Not found")
	ErrConflict          = errors.New("Already exists")
	ErrUnauthorized      = errors.New("Unauthorized")
	ErrForbidden         = errors.New("Forbidden")
	ErrValidation         = errors.New("Validation failed")
	ErrInvalidCredential = errors.New("Invalid credentail")
	ErrInvalidTransition = errors.New("Invalid transition")
)
