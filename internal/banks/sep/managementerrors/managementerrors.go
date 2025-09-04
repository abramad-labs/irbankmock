package managementerrors

import "errors"

var ErrEmptyName = errors.New("name can't be empty")
var ErrInvalidName = errors.New("name is not valid")

var ErrTokenNotFound = errors.New("token not found")
var ErrTokenExpired = errors.New("token expired")
var ErrTokenNoLongerAvailable = errors.New("token no longer available")

var ErrTransactionNotFound = errors.New("transaction not found")
