package http

import (
	"net/http"

	"github.com/pkg/errors"

	cError "github.com/lruggieri/fxnow/common/error"
)

func GetHttpStatusFromError(err error) int {
	if errors.Is(err, cError.Duplicated) {
		return http.StatusConflict
	} else if errors.Is(err, cError.NotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, cError.NotAuthenticated) {
		return http.StatusUnauthorized
	} else if errors.Is(err, cError.NotAuthorized) {
		return http.StatusForbidden
	} else if errors.Is(err, cError.InvalidParameter) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
