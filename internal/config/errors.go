package config

import "errors"

var (
	ErrDatabaseURLMissing = errors.New("URLDATABASE is not set")
	ErrJWTSecretMissing   = errors.New("JWT_SECRET is not set")
)
