package config

import "errors"

var ErrDatabaseURLMissing = errors.New("DATABASE_URL is not set")
