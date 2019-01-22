package opsgenie

import (
	"github.com/giantswarm/microerror"
)

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

var invalidTemplateError = &microerror.Error{
	Kind: "invalidTemplateError",
}

// IsInvalidTemplate asserts invalidTemplateError.
func IsInvalidTemplate(err error) bool {
	return microerror.Cause(err) == invalidTemplateError
}

var routingRuleDuplicationError = &microerror.Error{
	Kind: "routingRuleDuplicationError",
}

// IsRoutingRuleDuplication asserts routingRuleDuplication.
func IsRoutingRuleDuplication(err error) bool {
	return microerror.Cause(err) == routingRuleDuplicationError
}

var unexpectedResponseCodeError = &microerror.Error{
	Kind: "unexpectedResponseCodeError",
}

// IsUnexpectedResponseCode asserts unexpectedResponseCodeError.
func IsUnexpectedResponseCode(err error) bool {
	return microerror.Cause(err) == unexpectedResponseCodeError
}
