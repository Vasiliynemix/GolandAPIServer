package repo

import "errors"

const (
	ErrDuplicateKeyStr = "duplicate key value violates unique constraint"
)

var (
	ErrNotFound = errors.New("not found")
	ErrIsExists = errors.New("already exists")
)
