package apperr

import "errors"

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrRestaurantNotFound = errors.New("restaurant not found")
	ErrMenuNotFound       = errors.New("menu not found")
	ErrInvalidID          = errors.New("invalid id")
)
