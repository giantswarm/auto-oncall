package webhook

import (
	"github.com/giantswarm/microerror"
)

var decodeFailedError = &microerror.Error{
	Kind: "decodeFailedError",
}

// IsDecodeFailed asserts deleteFailedError.
func IsDecodeFailed(err error) bool {
	return microerror.Cause(err) == decodeFailedError
}
