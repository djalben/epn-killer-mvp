package environment

import "errors"

var (
	ErrPostgresDSNNotFound   = errors.New("postgres DSN not found in environment variables")
)
