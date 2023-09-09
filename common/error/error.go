package error

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicated       = errors.New("duplicated")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNotAuthorized    = errors.New("not authorized")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrTooManyRequests  = errors.New("too many requests")
)
