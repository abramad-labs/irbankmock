package usererror

import (
	"errors"
)

type UserError struct {
	err    error
	Status *int
}

func New(err error) error {
	return &UserError{
		err:    err,
		Status: nil,
	}
}

func NewBadRequest(err error) error {
	return NewWithStatus(err, 400)
}

func NewWithStatus(err error, status int) error {
	return &UserError{
		err:    err,
		Status: &status,
	}
}

func (e UserError) Error() string {
	return e.err.Error()
}

func (e UserError) Unwrap() error {
	return e.err
}

func (e UserError) Wrap(err error) error {
	var newError error
	if e.err == nil {
		newError = err
	} else {
		newError = errors.Join(e.err, err)
	}

	return &UserError{
		Status: e.Status,
		err:    newError,
	}
}

func (e UserError) Is(err error) bool {
	return errors.Is(e.Unwrap(), err)
}
