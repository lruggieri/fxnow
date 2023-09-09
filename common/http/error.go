package http

import (
	"net/http"

	"github.com/pkg/errors"

	cError "github.com/lruggieri/fxnow/common/error"
)

//nolint:revive
func GetHttpStatusFromError(err error) int {
	if errors.Is(err, cError.ErrDuplicated) {
		return http.StatusConflict
	} else if errors.Is(err, cError.ErrNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, cError.ErrNotAuthenticated) {
		return http.StatusUnauthorized
	} else if errors.Is(err, cError.ErrNotAuthorized) {
		return http.StatusForbidden
	} else if errors.Is(err, cError.ErrInvalidParameter) {
		return http.StatusBadRequest
	} else if errors.Is(err, cError.ErrTooManyRequests) {
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}
