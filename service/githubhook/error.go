package githubhook

import (
	"github.com/giantswarm/microerror"
)

var invalidHookError = &microerror.Error{
	Kind: "invalidHookError",
}

// IsInvalidHook asserts invalidHookError.
func IsInvalidHook(err error) bool {
	return microerror.Cause(err) == invalidHookError
}
