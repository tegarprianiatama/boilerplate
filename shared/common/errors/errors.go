package errors

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInvalidID              = errors.New("invalid user ID")
	ErrInvalidInput           = errors.New("invalid request payload")
	ErrUnauthorized           = errors.New("unauthorized access")
	ErrForbidden              = errors.New("forbidden access")
	ErrInternalServer         = errors.New("internal server error")
)
