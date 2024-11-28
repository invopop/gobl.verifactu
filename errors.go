package verifactu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/invopop/gobl.verifactu/internal/gateways"
)

// Main error types return by this package.
var (
	ErrValidation = newError("validation")
	ErrDuplicate  = newError("duplicate")
	ErrConnection = newError("connection")
	ErrInternal   = newError("internal")
)

// Error allows for structured responses to better handle errors upstream.
type Error struct {
	key     string
	code    string
	message string
	cause   error
}

func newError(key string) *Error {
	return &Error{key: key}
}

// newErrorFrom attempts to wrap the provided error into the Error type.
func newErrorFrom(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	} else if e, ok := err.(*gateways.Error); ok {
		return &Error{
			key:     e.Key(),
			code:    e.Code(),
			message: e.Message(),
			cause:   e,
		}
	}
	return &Error{
		key:     "internal",
		message: err.Error(),
		cause:   err,
	}
}

// Error produces a human readable error message.
func (e *Error) Error() string {
	out := []string{e.key}
	if e.code != "" {
		out = append(out, e.code)
	}
	if e.message != "" {
		out = append(out, e.message)
	}
	return strings.Join(out, ": ")
}

// Key returns the key for the error.
func (e *Error) Key() string {
	return e.key
}

// Message returns the human message for the error.
func (e *Error) Message() string {
	return e.message
}

// Code returns the code provided by the remote service.
func (e *Error) Code() string {
	return e.code
}

// Cause returns the undlying error that caused this error.
func (e *Error) Cause() error {
	return e.cause
}

// withMessage duplicates and adds the message to the error.
func (e *Error) withMessage(msg string, args ...any) *Error {
	e = e.clone()
	e.message = fmt.Sprintf(msg, args...)
	return e
}

func (e *Error) clone() *Error {
	ne := new(Error)
	*ne = *e
	return ne
}

// Is checks to see if the target error is the same as the current one
// or forms part of the chain.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.cause, target)
	}
	return e.key == t.key
}
