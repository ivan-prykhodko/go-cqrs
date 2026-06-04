package cqrs

import "errors"

var (
	ErrUnexpectedResultType = errors.New("handler returned unexpected result type")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvalidInputType     = errors.New("invalid input type")
	ErrHandlerNotFound      = errors.New("handler not found")
)
