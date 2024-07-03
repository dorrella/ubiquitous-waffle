package types

import (
	"errors"
)

var (
	ErrEmailValidation    = errors.New("invalid email")
	ErrCustomerValidation = errors.New("invalid customer")
	ErrDatabaseErr        = errors.New("db call failed")
	ErrUnexpectedResult   = errors.New("unexpected result")
)
