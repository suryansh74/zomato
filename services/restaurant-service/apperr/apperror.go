package apperr

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrUserNotFound   = errors.New("user not found")
)
