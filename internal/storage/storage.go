package storage

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
	ErrAppNotFound      = errors.New("application not found")
)
