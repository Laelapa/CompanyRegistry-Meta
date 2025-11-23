package domain

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource already exists")
	ErrBadCredentials = errors.New("invalid credentials")
	ErrBadRequest     = errors.New("bad request")
)
