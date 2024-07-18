package storage

import "errors"

var (
	ErrURLNotFound error = errors.New("url not found")
	ErrAliasExists error = errors.New("alias already exists")
)
