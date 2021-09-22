package nuwa

import (
	goerrors "errors"

	pkgerrors "github.com/pkg/errors"
)

// IsErr alias the errors.Is
var IsErr = goerrors.Is

// NewErr alias the errors.New
var NewErr = goerrors.New

// WrapfErr alias the Wrapf
var WrapfErr = pkgerrors.Wrapf

// Errorf alias the Errorf
var Errorf = pkgerrors.Errorf
