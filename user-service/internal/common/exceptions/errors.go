package exceptions

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrUnavailable = errors.New("user service unavailable")
var ErrInvalidUser = errors.New("user does not exist")
