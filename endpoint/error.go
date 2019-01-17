package endpoint

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var userNotFoundError = &microerror.Error{
	Kind: "userNotFoundError",
}

// IsUserNotFound asserts userNotFoundError.
func IsUserNotFound(err error) bool {
	return microerror.Cause(err) == userNotFoundError
}
