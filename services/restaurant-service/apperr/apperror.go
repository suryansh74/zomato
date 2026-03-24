package apperr

import "errors"

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrRestaurantNotFound = errors.New("restaurant not found")
)
