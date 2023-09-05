package error

import "errors"

var (
	NotFound         = errors.New("not found")
	Duplicated       = errors.New("duplicated")
	NotAuthenticated = errors.New("not authenticated")
	NotAuthorized    = errors.New("not authorized")
	InvalidParameter = errors.New("invalid parameter")
)
