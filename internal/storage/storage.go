package storage

import "errors"

var (
	ErrURLNotFound error = errors.New("url not found")
	ErrURLExists   error = errors.New("url already exists")
)
