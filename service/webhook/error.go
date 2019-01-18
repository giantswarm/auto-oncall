package webhook

import (
	"github.com/giantswarm/microerror"
)

var userNotFoundError = &microerror.Error{
	Kind: "userNotFoundError",
}

// IsUserNotFound asserts userNotFoundError.
func IsUserNotFound(err error) bool {
	return microerror.Cause(err) == userNotFoundError
}

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
