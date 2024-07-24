package services

import "errors"

var (
	ErrEnableToSave      = errors.New("can't save url")
	ErrAliasNotFound     = errors.New("alias not found")
	ErrFailedToFindAlias = errors.New("failed to find alias")
)
