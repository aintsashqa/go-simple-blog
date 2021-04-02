package errors

import (
	"errors"
)

var (
	ErrInternal error = errors.New("Internal server error")

	ErrInvalidTokenUserId     error = errors.New("Invalid token user id")
	ErrUnavailableRequestBody error = errors.New("Unavailable request body")
	ErrInvalidRequestBody     error = errors.New("Invalid request body")

	ErrInvalidAuthorizedUserID    error = errors.New("Invalid authorized user id")
	ErrEmptyAuthorizationHeader   error = errors.New("Header `Authorization` could not be empty")
	ErrInvalidAuthorizationHeader error = errors.New("Invalid `Authorization` header")
	ErrAuthenticationFailed       error = errors.New("Authentication failed")
)
