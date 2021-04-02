package errors

import (
	"errors"
)

var (
	ErrUserNotFound error = errors.New("User not found is database")
	ErrPostNotFound error = errors.New("Post not found in database")
)
