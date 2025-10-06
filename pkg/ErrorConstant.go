package pkg

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmptyPassword     = errors.New("password cannot be empty")
)
