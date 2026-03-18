package registry

import "errors"

var (
	ErrDuplicateKey = errors.New("registry: duplicate key")
	ErrNotFound     = errors.New("registry: key not found")
)
