package models

import (
	"errors"
	"net/http"
)

var (
	ErrEmailExists  error = errors.New("models: users with such email already exists")
	ErrNotFound     error = errors.New("models: couldn't find an entity")
	ErrBadRequest   error = errors.New("models: invalid request params")
	ErrUnauthorized error = errors.New("models: insufficient level of access")
)

func HttpErrorCode(err error) int {
	switch {
	case errors.Is(err, ErrEmailExists):
		return http.StatusConflict
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrUnauthorized):
		return http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
