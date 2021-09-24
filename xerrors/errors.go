package xerrors

import (
	errors "errors"

	pkgerrors "github.com/pkg/errors"
)

// Is alias the errors.Is
var Is = errors.Is

// New alias the errors.New
var New = errors.New

// Wrapf alias the Wrapf
var Wrapf = pkgerrors.Wrapf

// Errorf alias the Errorf
var Errorf = pkgerrors.Errorf

var (
	// ErrNotImplement defines the function is not impelement
	ErrNotImplement = errors.New("not implement")

	// ErrNotFound defines the object is not found
	ErrNotFound = errors.New("not found")
)

// WrapNotFound return the wraped not found error
func WrapNotFound(format string, args ...interface{}) error {
	return Wrapf(ErrNotFound, format, args...)
}
