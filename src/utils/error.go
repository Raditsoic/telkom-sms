package utils

import "errors"

var ErrUsernameExists = errors.New("username already exists")

var ErrInvalidCredentials = errors.New("invalid username or password")

var ErrInvalidID = errors.New("invalid ID")