package models

import "errors"

var (
	ErrEmailExists error = errors.New("models: users with such email already exists")
	ErrNotFound    error = errors.New("models: couldn't find an entity")
)
