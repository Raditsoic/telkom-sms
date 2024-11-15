package utils

import "errors"

var ErrUsernameExists = errors.New("username already exists")

var ErrInvalidCredentials = errors.New("invalid username or password")

var ErrInvalidID = errors.New("invalid ID")

var ErrTransactionType = errors.New("invalid transaction type")

var ErrTransactionNotFound = errors.New("transaction not found")

var ErrItemNotFound = errors.New("item not found")